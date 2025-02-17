package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient(kubeconfig string) (*Client, error) {
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) CreateSecret(ctx context.Context, data *SecretData) error {
	_, err := c.GetSecret(ctx, data.Namespace, data.Name)
	if err == nil {
		return fmt.Errorf("secret %s already exists in namespace %s", data.Name, data.Namespace)
	}
	if !errors.IsNotFound(err) {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Type: corev1.SecretType(data.Type),
		Data: makeSecretData(data.Data),
	}

	_, err = c.clientset.CoreV1().Secrets(data.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating secret: %w", err)
	}

	return nil
}

func (c *Client) UpdateSecret(ctx context.Context, data *SecretData) error {
	existing, err := c.GetSecret(ctx, data.Namespace, data.Name)
	if err != nil {
		return fmt.Errorf("error getting existing secret: %w", err)
	}

	existing.Data = makeSecretData(data.Data)
	if data.Type != "" {
		existing.Type = corev1.SecretType(data.Type)
	}

	_, err = c.clientset.CoreV1().Secrets(data.Namespace).Update(ctx, existing, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating secret: %w", err)
	}

	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, namespace, name string) error {
	err := c.clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("secret %s not found in namespace %s", name, namespace)
		}
		return fmt.Errorf("error deleting secret: %w", err)
	}

	return nil
}

func (c *Client) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("secret %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("error getting secret: %w", err)
	}

	return secret, nil
}

func (c *Client) ListSecrets(ctx context.Context, namespace string) ([]corev1.Secret, error) {
	secretList, err := c.clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing secrets: %w", err)
	}

	return secretList.Items, nil
}

func makeSecretData(data map[string]string) map[string][]byte {
	secretData := make(map[string][]byte)
	for key, value := range data {
		secretData[key] = []byte(value)
	}
	return secretData
}

func (c *Client) GetSecretString(ctx context.Context, namespace, name, key string) (string, error) {
	secret, err := c.GetSecret(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	value, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s", key, name)
	}

	return string(value), nil
}

func (c *Client) Exists(ctx context.Context, namespace, name string) (bool, error) {
	_, err := c.GetSecret(ctx, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
