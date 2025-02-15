package cmd

import (
	"context"
	"fmt"

	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deleta um secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("error creating k8s client: %w", err)
		}

		if err := client.DeleteSecret(context.Background(), namespace, secretName); err != nil {
			return fmt.Errorf("error deleting secret: %w", err)
		}

		fmt.Printf("Secret %s deletado com sucesso do namespace %s\n", secretName, namespace)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&secretName, "name", "", "nome do secret")
	deleteCmd.MarkFlagRequired("name")
}
