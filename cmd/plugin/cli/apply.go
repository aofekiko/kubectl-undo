package cli

import (
	"github.com/aofekiko/kubectl-undo/pkg/logger"
	"github.com/spf13/cobra"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply TYPE NAME VERSION",
	Short: "Apply an older version of a resource",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
		log.Info("implement apply")
	},
}
