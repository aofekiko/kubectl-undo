package request

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aofekiko/kubectl-undo/pkg/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	log = logger.NewLogger()
)

func DiscoverGroupVersion(discoveryClient discovery.DiscoveryClient, ResourceKind string) (*metav1.APIResource, error) {
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

func GetStaleResource(configFlags *genericclioptions.ConfigFlags, discoveryClient *discovery.DiscoveryClient, clientSet *kubernetes.Clientset, ResourceKind string, ResourceName string, ResourceVersion string) (*unstructured.Unstructured, error) {
	//TODO: correct logging
	//TODO: Error handling
	//TODO: make a controller to have the clients present instead of as parameters

	//discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)

	apiVersion, err := DiscoverGroupVersion(*discoveryClient, ResourceKind)
	if err != nil {
		log.Info("Failed to list API Resources")
		//log.Info(fmt.Sprintf("Failed to list API Resources: %v\n", err))
	}

	objects := &unstructured.UnstructuredList{}
	request := clientSet.RESTClient().Get().Prefix(fmt.Sprintf("/api/%s/%s", apiVersion.Group, apiVersion.Version)).Resource(apiVersion.Name).NamespaceIfScoped(*configFlags.Namespace, apiVersion.Namespaced).Param("resourceVersion", ResourceVersion).Param("resourceVersionMatch", "Exact")
	URL := request.URL()
	log.Info(URL.String())
	err = request.Do(context.TODO()).Into(objects)
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

func GetCurrentResource(discoveryClient *discovery.DiscoveryClient, dynamicClient *dynamic.DynamicClient, resourceKind string, resourceName string, namespace string) (*unstructured.Unstructured, error) {
	resource, err := DiscoverGroupVersion(*discoveryClient, resourceKind)
	if err != nil {
		log.Info("Failed to list API Resources")
		//log.Info(fmt.Sprintf("Failed to list API Resources: %v\n", err))
		return nil, err
	}
	unstructured, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    resource.Group,
		Version:  resource.Version,
		Resource: resource.Name,
	}).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})
	if err != nil {
		log.Info("Failed to get resource")
		//log.Info(fmt.Sprintf("Failed to get resource: %v\n", err))
		return nil, err
	}
	return unstructured, nil
}

func GetMostRecentStaleResource(configFlags *genericclioptions.ConfigFlags, dynamicClient *dynamic.DynamicClient, discoveryClient *discovery.DiscoveryClient, clientSet *kubernetes.Clientset, ResourceKind string, ResourceName string, ResourceVersion string) (*unstructured.Unstructured, error) {
	unstructured, err := GetCurrentResource(discoveryClient, dynamicClient, ResourceKind, ResourceName, *configFlags.Namespace)
	if err != nil {
		log.Info("Failed to get object")
		//log.Info(fmt.Sprintf("Failed to get object: %v\n", err))
		return nil, err
	}
	ResourceVersion = unstructured.GetResourceVersion()
	ResourceVersionInt, err := strconv.Atoi(ResourceVersion)
	if err != nil {
		log.Info("Failed to parse object's resourceVersion")
		//log.Info(fmt.Sprintf("Failed to parse object's resourceVersion: %v\n", err))
		return nil, err
	}
	ResourceVersionInt--
	ResourceVersion = strconv.Itoa(ResourceVersionInt)
	resource, err := GetStaleResource(configFlags, discoveryClient, clientSet, ResourceKind, ResourceName, ResourceVersion)
	return resource, err
}
