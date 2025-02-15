package cmd

import (
	"context"
	"fmt"

	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista secrets em um namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("error creating k8s client: %w", err)
		}

		secrets, err := client.ListSecrets(context.Background(), namespace)
		if err != nil {
			return fmt.Errorf("error listing secrets: %w", err)
		}

		fmt.Printf("Secrets no namespace %s:\n", namespace)
		for _, secret := range secrets {
			fmt.Printf("- %s (Tipo: %s)\n", secret.Name, secret.Type)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
