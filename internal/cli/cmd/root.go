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
	Short: "Kubernetes Secret Manager - Manage Kubernetes Secrets",
}

func Execute(cfg *config.Config) error {
	// Store config in package-level variable
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.Root().SetContext(cmd.Context())
	}
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file path")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "kubernetes namespace")
}

func initConfig() {
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	}
}
