package cli

import (
	//"errors"
	"fmt"
	"os"
	"strings"

	//"time"

	//"github.com/aofekiko/kubectl-undo/pkg/logger"
	//"github.com/aofekiko/kubectl-undo/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//"github.com/tj/go-spin"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

var (
	undoShort = `Undo recent resource changes`
	undoLong  = `
	Undo resource changes and revert them to a previous version.
	
	As long as etcd has not been compacted (https://etcd.io/docs/latest/op-guide/maintenance/#auto-compaction)`
	undoExample = `
	kubectl undo get configmap myconfigmap 1

	kubectl apply get configmap myconfigmap 1
	`
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "undo",
		Short:         undoShort,
		Long:          undoLong,
		Example:       undoExample,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(
		ApplyCmd,
		GetCmd,
	)

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.PersistentFlags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
