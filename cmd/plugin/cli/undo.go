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
	"github.com/aofekiko/kubectl-undo/pkg/logger"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	ClientSet             *kubernetes.Clientset
	DiscoveryClient       *discovery.DiscoveryClient
	DynamicClient         *dynamic.DynamicClient
	OutputFlag            string
)

var (
	undoShort = `Undo recent resource changes`
	undoLong  = `
	Undo resource changes and revert them to a previous version.

	the "-o yaml" flag is implied
	
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

	log := logger.NewLogger()

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	KubernetesConfigFlags.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().StringVarP(&OutputFlag, "output", "o", "yaml", "Output format. One of: (json, yaml)")

	config, err := KubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		log.Info("failed to read kubeconfig")
		//log.Info(fmt.Sprintf("failed to read kubeconfig: %w", err))
	}
	DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(config)

	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Info("failed to create clientset")
		//log.Info(fmt.Sprintf("Error creating dynamic client: %v\n", err))
	}

	DynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Info("failed to create dynamicclient")
	}

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
