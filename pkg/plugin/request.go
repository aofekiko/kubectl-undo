package request

import (
	"fmt"
	"os"
	"strings"

	"github.com/aofekiko/kubectl-undo/pkg/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func getApiResources(discoveryClient discovery.DiscoveryClient, ResourceKind string) (*metav1.APIResource, error) {
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
				resource.Group = gv.Group
				resource.Version = gv.Version
				return &resource, nil
			} else {
				for _, shortName := range resource.ShortNames {
					if shortName == lowercaseKind {
						resource.Group = gv.Group
						resource.Version = gv.Version
						return &resource, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func BuildRequest(configFlags *genericclioptions.ConfigFlags, ResourceKind string, ResourceName string, ResourceVersion string) (*unstructured.Unstructured, error) {
	log := logger.NewLogger()
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		log.Info("failed to read kubeconfig")
		//log.Info(fmt.Sprintf("failed to read kubeconfig: %w", err))
		return nil, err
	}

	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)

	apiVersion, err := getApiResources(*discoveryClient, ResourceKind)
	if err != nil {
		log.Info("Failed to list API Resources")
		//log.Info(fmt.Sprintf("Failed to list API Resources: %v\n", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Info("failed to create clientset")
		//log.Info(fmt.Sprintf("Error creating dynamic client: %v\n", err))
	}

	object := &unstructured.Unstructured{}
	err = clientset.RESTClient().Get().Prefix(fmt.Sprintf("/api/%s/%s", apiVersion.Group, apiVersion.Version)).Resource(apiVersion.Name).NamespaceIfScoped(*configFlags.Namespace, apiVersion.Namespaced).SetHeader("resourceVersion", ResourceVersion).SetHeader("resourceVersionMatch", "Exact").Do().Into(object)
	if err != nil {
		log.Info("failed to get resource")
		//log.Info(fmt.Sprintf("failed to get resource: %v\n", err))
	}
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	err = serializer.Encode(object, os.Stdout)
	if err != nil {
		log.Info(fmt.Sprintf("failed to encode objects: %v\n", err))
	}

	return nil, nil
}
