package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type SecretData struct {
	Name      string            `json:"name" validate:"required"`
	Namespace string            `json:"namespace" validate:"required"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data" validate:"required"`
}

type SecretManager interface {
	CreateSecret(ctx context.Context, data *SecretData) error

	UpdateSecret(ctx context.Context, data *SecretData) error

	DeleteSecret(ctx context.Context, namespace, name string) error

	GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error)

	ListSecrets(ctx context.Context, namespace string) ([]corev1.Secret, error)
}
