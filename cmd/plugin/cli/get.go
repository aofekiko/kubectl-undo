package cli

import (
	"fmt"
	"os"

	"github.com/aofekiko/kubectl-undo/pkg/logger"
	request "github.com/aofekiko/kubectl-undo/pkg/plugin"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var GetCmd = &cobra.Command{
	Use:   "get TYPE NAME [RESOURCEVERSION]",
	Short: "Get an older version of a resource, no resource version will get the most recent stale resource",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
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
		var serializer *json.Serializer
		if OutputFlag == "yaml" {
			serializer = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		} else if OutputFlag == "json" {
			serializer = json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
		}
		err := serializer.Encode(resource, os.Stdout)
		if err != nil {
			log.Info(fmt.Sprintf("failed to encode objects: %v\n", err))
		}
	},
}
