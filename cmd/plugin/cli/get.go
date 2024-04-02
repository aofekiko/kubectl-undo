package cli

import (
	"github.com/aofekiko/kubectl-undo/pkg/logger"
	request "github.com/aofekiko/kubectl-undo/pkg/plugin"
	"github.com/spf13/cobra"
)

var GetCmd = &cobra.Command{
	Use:   "get TYPE NAME VERSION",
	Short: "Get an older version of a resource",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
		log.Info("starting request")
		ResourceType := args[0]
		ResourceName := args[1]
		ResourceVersion := args[2]
		_, err := request.BuildRequest(KubernetesConfigFlags, ResourceType, ResourceName, ResourceVersion)
		if err != nil {
			log.Error(err)
		}
	},
}
