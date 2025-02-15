package cmd

import (
	"github.com/mpalu/k8s-secrets-manager/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	kubeconfig string
	namespace  string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-secrets-manager",
	Short: "Kubernetes Secret Manager - Manager Kubernetes Secrets",
	Long: `Kubernetes Secret Manager is a tool to create, update,
delete and list Kubernetes secrets. Can be used as a CLI tool or as a HTTP server.`,
}

func Execute(cfg *config.Config) error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-secrets-manager.yaml)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "kubernetes namespace")
}

func initConfig() {
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	}
}
