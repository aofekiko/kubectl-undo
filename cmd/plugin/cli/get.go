package cli

import (
	"fmt"
	"os"

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
		ResourceType := args[0]
		ResourceName := args[1]
		ResourceVersion := args[2]
		resource, err := request.GetStaleResource(KubernetesConfigFlags, ResourceType, ResourceName, ResourceVersion)
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
