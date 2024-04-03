package cli

import (
	printer "github.com/aofekiko/kubectl-undo/pkg/printer"
	request "github.com/aofekiko/kubectl-undo/pkg/request"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/cobra"
)

var DiffCmd = &cobra.Command{
	Use:   "diff TYPE NAME [RESOURCEVERSION]",
	Short: "Compare the currrent resource to an older version of the resource ",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		ResourceKind := args[0]
		ResourceName := args[1]
		var oldResource string
		var ResourceVersion string
		if len(args) == 2 { //TODO: change the trigger for this feature to a negative number instead of 0
			res, err := request.GetMostRecentStaleResource(KubernetesConfigFlags, DynamicClient, DiscoveryClient, ClientSet, ResourceKind, ResourceName, ResourceVersion)
			if err != nil {
				log.Info("failed to get the most recent stale resource")
			}
			oldResource, err = printer.PrintUnstructured(DiffOutputFlag, res)
			if err != nil {
				log.Info("Failed to stringify resource")
			}
		} else {
			ResourceVersion = args[2]
			res, err := request.GetStaleResource(KubernetesConfigFlags, DiscoveryClient, ClientSet, ResourceKind, ResourceName, ResourceVersion)
			if err != nil {
				log.Info("Failed to get stale resource")
			}
			oldResource, err = printer.PrintUnstructured(DiffOutputFlag, res)
			if err != nil {
				log.Info("Failed to stringify resource")
			}
		}
		currentResource, err := request.GetCurrentResource(DiscoveryClient, DynamicClient, ResourceKind, ResourceName, *KubernetesConfigFlags.Namespace)
		if err != nil {
			log.Info("Failed to get current resource")
		}
		currentResourceString, err := printer.PrintUnstructured(DiffOutputFlag, currentResource)
		if err != nil {
			log.Info("Failed to stringify resource")
		}
		if err != nil {
			log.Info("Failed to get current resource")
		}

		print(diff.Diff(currentResourceString, oldResource))
		//TODO: mark the + & - with () so not to confuse them with yaml arrays
	},
}
