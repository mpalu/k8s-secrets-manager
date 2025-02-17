package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockClient implements k8s.Client interface for testing
type mockClient struct {
	secrets map[string]*k8s.SecretData
}

func newMockClient() *mockClient {
	return &mockClient{
		secrets: make(map[string]*k8s.SecretData),
	}
}

func (m *mockClient) CreateSecret(ctx context.Context, data *k8s.SecretData) error {
	key := data.Namespace + "/" + data.Name
	if _, exists := m.secrets[key]; exists {
		return &k8s.ValidationError{
			Field:   "name",
			Message: "secret already exists",
		}
	}
	m.secrets[key] = data
	return nil
}

func (m *mockClient) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	key := namespace + "/" + name
	if secret, exists := m.secrets[key]; exists {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
			Data: map[string][]byte{
				"key1": []byte("value1"),
			},
		}, nil
	}
	return nil, &k8s.NotFoundError{
		Resource: "secret",
		Name:     name,
	}
}

func (m *mockClient) DeleteSecret(ctx context.Context, namespace, name string) error {
	key := namespace + "/" + name
	if _, exists := m.secrets[key]; !exists {
		return &k8s.ValidationError{
			Field:   "name",
			Message: "secret not found",
		}
	}
	delete(m.secrets, key)
	return nil
}

func (m *mockClient) ListSecrets(ctx context.Context, namespace string) ([]corev1.Secret, error) {
	secrets := []corev1.Secret{}
	for _, secret := range m.secrets {
		if secret.Namespace == namespace {
			secrets = append(secrets, corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secret.Name,
					Namespace: secret.Namespace,
				},
				Data: map[string][]byte{
					"key1": []byte("value1"),
				},
			})
		}
	}
	return secrets, nil
}

func (m *mockClient) UpdateSecret(ctx context.Context, data *k8s.SecretData) error {
	key := data.Namespace + "/" + data.Name
	if _, exists := m.secrets[key]; !exists {
		return &k8s.ValidationError{
			Field:   "name",
			Message: "secret not found",
		}
	}
	m.secrets[key] = data
	return nil
}

func TestCreateSecret(t *testing.T) {
	mockClient := newMockClient()
	handler := NewHandler(mockClient)

	tests := []struct {
		name           string
		secretData     k8s.SecretData
		expectedStatus int
	}{
		{
			name: "valid secret",
			secretData: k8s.SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Type:      "Opaque",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid secret - missing name",
			secretData: k8s.SecretData{
				Namespace: "default",
				Type:      "Opaque",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.secretData)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/secrets", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler.CreateSecret(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetSecret(t *testing.T) {
	mockClient := newMockClient()
	handler := NewHandler(mockClient)

	// Create a test secret
	testSecret := &k8s.SecretData{
		Name:      "test-secret",
		Namespace: "default",
		Type:      "Opaque",
		Data: map[string]string{
			"key1": "value1",
		},
	}
	mockClient.CreateSecret(context.Background(), testSecret)

	tests := []struct {
		name           string
		secretName     string
		namespace      string
		expectedStatus int
	}{
		{
			name:           "existing secret",
			secretName:     "test-secret",
			namespace:      "default",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existing secret",
			secretName:     "non-existing",
			namespace:      "default",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/secrets/"+tt.secretName+"?namespace="+tt.namespace, nil)
			rr := httptest.NewRecorder()

			// Setup router to handle path variables
			router := mux.NewRouter()
			router.HandleFunc("/api/v1/secrets/{name}", handler.GetSecret)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	mockClient := newMockClient()
	handler := NewHandler(mockClient)

	// Create initial secret
	initialSecret := &k8s.SecretData{
		Name:      "test-secret",
		Namespace: "default",
		Type:      "Opaque",
		Data: map[string]string{
			"key1": "value1",
		},
	}
	mockClient.CreateSecret(context.Background(), initialSecret)

	tests := []struct {
		name           string
		secretName     string
		secretData     k8s.SecretData
		expectedStatus int
	}{
		{
			name:       "successful update",
			secretName: "test-secret",
			secretData: k8s.SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Data: map[string]string{
					"key1": "updated-value",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "non-existent secret",
			secretName: "non-existent",
			secretData: k8s.SecretData{
				Name:      "non-existent",
				Namespace: "default",
				Type:      "Opaque",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.secretData)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/secrets/"+tt.secretName, bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/v1/secrets/{name}", handler.UpdateSecret)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
