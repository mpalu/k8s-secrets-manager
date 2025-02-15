package validator

import (
	"fmt"

	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
)

func ValidateSecretData(data *k8s.SecretData) error {
	if data.Name == "" {
		return &k8s.ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}

	if data.Namespace == "" {
		return &k8s.ValidationError{
			Field:   "namespace",
			Message: "namespace is required",
		}
	}

	if len(data.Data) == 0 {
		return &k8s.ValidationError{
			Field:   "data",
			Message: "at least one data entry is required",
		}
	}

	for key, value := range data.Data {
		if key == "" {
			return &k8s.ValidationError{
				Field:   "data",
				Message: "empty key is not allowed",
			}
		}

		if !isValidKey(key) {
			return &k8s.ValidationError{
				Field:   "data",
				Message: fmt.Sprintf("invalid key format: %s", key),
			}
		}
	}

	return nil
}

func isValidKey(key string) bool {
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
	return true
}
