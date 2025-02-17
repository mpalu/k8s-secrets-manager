package validator

import (
	"testing"

	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
)

func TestValidateSecretData(t *testing.T) {
	tests := []struct {
		name     string
		data     *k8s.SecretData
		wantErr  bool
		errField string
	}{
		{
			name: "valid secret data",
			data: &k8s.SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Type:      "Opaque",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			data: &k8s.SecretData{
				Namespace: "default",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			wantErr:  true,
			errField: "name",
		},
		{
			name: "missing namespace",
			data: &k8s.SecretData{
				Name: "test-secret",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			wantErr:  true,
			errField: "namespace",
		},
		{
			name: "empty data",
			data: &k8s.SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Data:      map[string]string{},
			},
			wantErr:  true,
			errField: "data",
		},
		{
			name: "empty key in data",
			data: &k8s.SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Data: map[string]string{
					"": "value1",
				},
			},
			wantErr:  true,
			errField: "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecretData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if validErr, ok := err.(*k8s.ValidationError); ok {
					if validErr.Field != tt.errField {
						t.Errorf("ValidateSecretData() error field = %v, want %v", validErr.Field, tt.errField)
					}
				}
			}
		})
	}
}
