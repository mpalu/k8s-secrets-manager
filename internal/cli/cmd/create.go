package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	"github.com/spf13/cobra"
)

var (
	secretName string
	secretType string
	secretData string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient(kubeconfig)
		if err != nil {
			return fmt.Errorf("error creating k8s client: %w", err)
		}

		data := make(map[string]string)
		pairs := strings.Split(secretData, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}

		secret := &k8s.SecretData{
			Name:      secretName,
			Namespace: namespace,
			Type:      secretType,
			Data:      data,
		}

		if err := client.CreateSecret(context.Background(), secret); err != nil {
			return fmt.Errorf("error creating secret: %w", err)
		}

		fmt.Printf("Secret %s successfully created in namespace %s\n", secretName, namespace)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&secretName, "name", "", "secret name")
	createCmd.Flags().StringVar(&secretType, "type", "Opaque", "secret type")
	createCmd.Flags().StringVar(&secretData, "data", "", "secret data (format: key1=value1,key2=value2)")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("data")
}
