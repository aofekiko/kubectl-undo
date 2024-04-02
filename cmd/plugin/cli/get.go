package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aofekiko/kubectl-undo/pkg/logger"
	request "github.com/aofekiko/kubectl-undo/pkg/plugin"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var GetCmd = &cobra.Command{
	Use:   "get TYPE NAME VERSION",
	Short: "Get an older version of a resource",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
		ResourceKind := args[0]
		ResourceName := args[1]
		ResourceVersion := args[2]
		if i, err := strconv.Atoi(ResourceVersion); err == nil && i == 0 { //TODO: change the trigger for this feature to a negative number instead of 0
			unstructured, err := request.GetCurrentResource(DiscoveryClient, DynamicClient, ResourceVersion, ResourceKind, ResourceName, *KubernetesConfigFlags.Namespace)
			ResourceVersion = unstructured.GetResourceVersion()
			ResourceVersionInt, err := strconv.Atoi(ResourceVersion)
			if err != nil {
				log.Info("Failed to parse object's resourceVersion")
				//log.Info(fmt.Sprintf("Failed to parse object's resourceVersion: %v\n", err))
			}
			ResourceVersionInt--
			ResourceVersion = strconv.Itoa(ResourceVersionInt)
		}
		resource, err := request.GetStaleResource(KubernetesConfigFlags, DiscoveryClient, ClientSet, ResourceKind, ResourceName, ResourceVersion)
		if err != nil {
			log.Error(err)
		}
		serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		err = serializer.Encode(resource, os.Stdout)
		if err != nil {
			log.Info(fmt.Sprintf("failed to encode objects: %v\n", err))
		}

	},
}
