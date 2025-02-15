package cmd

import (
	"fmt"

	"github.com/mpalu/k8s-secrets-manager/internal/api/server"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	"github.com/spf13/cobra"
)

var (
	port string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Inicia o servidor HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("error creating k8s client: %w", err)
		}

		srv := server.New(client)
		return srv.Run(":" + port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "HTTP server port")
}
