package request

import (
	"fmt"
	"strings"

	"github.com/aofekiko/kubectl-undo/pkg/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
)

func discoverGroupVersion(discoveryClient discovery.DiscoveryClient, ResourceKind string) (*metav1.APIResource, error) {
	lowercaseKind := strings.ToLower(ResourceKind)
	resources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	for _, ResourceGroupVersion := range resources {
		gv, err := schema.ParseGroupVersion(ResourceGroupVersion.GroupVersion)
		if err != nil {
			return nil, err
		}
		for _, resource := range ResourceGroupVersion.APIResources {
			if resource.Kind == lowercaseKind || resource.Name == lowercaseKind || resource.SingularName == lowercaseKind {
				return buildAPIResource(resource, gv)
			} else {
				for _, shortName := range resource.ShortNames {
					if shortName == lowercaseKind {
						return buildAPIResource(resource, gv)
					}
				}
			}
		}
	}
	return nil, nil
}

func buildAPIResource(resource metav1.APIResource, gv schema.GroupVersion) (*metav1.APIResource, error) {
	resource.Group = gv.Group
	resource.Version = gv.Version
	return &resource, nil
}

func GetStaleResource(configFlags *genericclioptions.ConfigFlags, ResourceKind string, ResourceName string, ResourceVersion string) (*unstructured.Unstructured, error) {
	log := logger.NewLogger()
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		log.Info("failed to read kubeconfig")
		//log.Info(fmt.Sprintf("failed to read kubeconfig: %w", err))
		return nil, err
	}

	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)

	apiVersion, err := discoverGroupVersion(*discoveryClient, ResourceKind)
	if err != nil {
		log.Info("Failed to list API Resources")
		//log.Info(fmt.Sprintf("Failed to list API Resources: %v\n", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Info("failed to create clientset")
		//log.Info(fmt.Sprintf("Error creating dynamic client: %v\n", err))
	}

	objects := &unstructured.UnstructuredList{}
	request := clientset.RESTClient().Get().Prefix(fmt.Sprintf("/api/%s/%s", apiVersion.Group, apiVersion.Version)).Resource(apiVersion.Name).NamespaceIfScoped(*configFlags.Namespace, apiVersion.Namespaced).Param("resourceVersion", ResourceVersion).Param("resourceVersionMatch", "Exact")
	URL := request.URL()
	log.Info(URL.String())
	err = request.Do().Into(objects)
	if err != nil {
		log.Info("failed to get resource, the ResourceVersion may have been compacted")
		//log.Info(fmt.Sprintf("failed to get resource: %v\n", err))
	}
	for _, resource := range objects.Items {
		if resource.GetName() == ResourceName {
			return &resource, nil
		}
	}
	log.Info("Did not find a resource with the matching name")
	return nil, nil
}
