package cli

import (
	"context"

	printer "github.com/aofekiko/kubectl-undo/pkg/printer"
	request "github.com/aofekiko/kubectl-undo/pkg/request"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply TYPE NAME [RESOURCEVERSION]",
	Short: "Apply an older version of a resource, no resource version will get the most recent stale resource",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		ResourceKind := args[0]
		ResourceName := args[1]
		var resource *unstructured.Unstructured
		var ResourceVersion string
		if len(args) == 2 { //TODO: change the trigger for this feature to a negative number instead of 0
			res, err := request.GetMostRecentStaleResource(KubernetesConfigFlags, DynamicClient, DiscoveryClient, ClientSet, ResourceKind, ResourceName, ResourceVersion)
			if err != nil {
				log.Info("failed to get the most recent stale resource")
			}
			resource = res
		} else {
			ResourceVersion = args[2]
			res, err := request.GetStaleResource(KubernetesConfigFlags, DiscoveryClient, ClientSet, ResourceKind, ResourceName, ResourceVersion)
			if err != nil {
				log.Info("Failed to get stale resource")
			}
			resource = res
		}

		gvr, err := request.KindToGVR(*DiscoveryClient, ResourceKind)
		if err != nil {
			log.Info("Failed to find the api resource by kind")
		}

		resource.SetResourceVersion("")
		resource.SetManagedFields(nil)
		metadata, found, err := unstructured.NestedMap(resource.Object, "metadata")
		if err != nil {
			log.Info("Failed to fetch object metadata")
		}
		if found {
			metadata["creationTimestamp"] = nil
			if err := unstructured.SetNestedMap(resource.Object, metadata, "metadata"); err != nil {
				log.Info("Failed to wipe creationTimestamp off applied resource")
			}
		}

		newResource, err := DynamicClient.Resource(*gvr).Namespace(resource.GetNamespace()).Apply(context.TODO(), resource.GetName(), resource, v1.ApplyOptions{FieldManager: fieldManager, Force: ForceFlag})
		//_, err = DynamicClient.Resource(*gvr).Namespace(resource.GetNamespace()).Update(context.TODO(), resource, v1.UpdateOptions{FieldManager: fieldManager})
		if err != nil {
			log.Info("failed to apply resource")
		}
		if OutputFlag != "none" {
			printer.PrintUnstructured(OutputFlag, newResource)
		}
	},
}
