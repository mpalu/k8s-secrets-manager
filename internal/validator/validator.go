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

	for key := range data.Data {
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
	// DNS subdomain name validation as per Kubernetes naming conventions
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
	if len(key) > 253 {
		return false
	}

	// Must consist of lowercase alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character
	validKey := true
	for i, char := range key {
		isAlphanumeric := (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')
		isDash := char == '-'
		isDot := char == '.'

		// First and last characters must be alphanumeric
		if i == 0 || i == len(key)-1 {
			validKey = validKey && isAlphanumeric
		} else {
			validKey = validKey && (isAlphanumeric || isDash || isDot)
		}
	}
	return validKey
}
