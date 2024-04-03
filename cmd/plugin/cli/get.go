package cli

import (
	printer "github.com/aofekiko/kubectl-undo/pkg/printer"
	request "github.com/aofekiko/kubectl-undo/pkg/request"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var GetCmd = &cobra.Command{
	Use:   "get TYPE NAME [RESOURCEVERSION]",
	Short: "Get an older version of a resource, no resource version will get the most recent stale resource",
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
				log.Error(err)
			}
			resource = res
		}
		printer.PrintUnstructured(OutputFlag, resource)
	},
}
