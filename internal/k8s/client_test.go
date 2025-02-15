package k8s

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClient_CreateSecret(t *testing.T) {
	tests := []struct {
		name        string
		secretData  *SecretData
		setupFunc   func(*fake.Clientset)
		wantErr     bool
		errorString string
	}{
		{
			name: "successful creation",
			secretData: &SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Type:      string(corev1.SecretTypeOpaque),
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			setupFunc: func(clientset *fake.Clientset) {},
			wantErr:   false,
		},
		{
			name: "secret already exists",
			secretData: &SecretData{
				Name:      "existing-secret",
				Namespace: "default",
				Type:      string(corev1.SecretTypeOpaque),
				Data: map[string]string{
					"key1": "value1",
				},
			},
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "existing-secret",
						Namespace: "default",
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantErr:     true,
			errorString: "secret existing-secret already exists in namespace default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			err := client.CreateSecret(context.TODO(), tt.secretData)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("CreateSecret() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr {
				secret, err := client.GetSecret(context.TODO(), tt.secretData.Namespace, tt.secretData.Name)
				if err != nil {
					t.Errorf("Failed to get created secret: %v", err)
					return
				}

				for k, v := range tt.secretData.Data {
					if string(secret.Data[k]) != v {
						t.Errorf("Secret data mismatch for key %s: got %s, want %s", k, secret.Data[k], v)
					}
				}
			}
		})
	}
}

func TestClient_UpdateSecret(t *testing.T) {
	tests := []struct {
		name        string
		secretData  *SecretData
		setupFunc   func(*fake.Clientset)
		wantErr     bool
		errorString string
	}{
		{
			name: "successful update",
			secretData: &SecretData{
				Name:      "test-secret",
				Namespace: "default",
				Type:      string(corev1.SecretTypeOpaque),
				Data: map[string]string{
					"key1": "new-value1",
				},
			},
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"key1": []byte("old-value1"),
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantErr: false,
		},
		{
			name: "secret not found",
			secretData: &SecretData{
				Name:      "non-existing-secret",
				Namespace: "default",
				Data: map[string]string{
					"key1": "value1",
				},
			},
			setupFunc:   func(clientset *fake.Clientset) {},
			wantErr:     true,
			errorString: "error getting existing secret: secret non-existing-secret not found in namespace default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			err := client.UpdateSecret(context.TODO(), tt.secretData)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("UpdateSecret() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr {
				secret, err := client.GetSecret(context.TODO(), tt.secretData.Namespace, tt.secretData.Name)
				if err != nil {
					t.Errorf("Failed to get updated secret: %v", err)
					return
				}

				for k, v := range tt.secretData.Data {
					if string(secret.Data[k]) != v {
						t.Errorf("Secret data mismatch for key %s: got %s, want %s", k, secret.Data[k], v)
					}
				}
			}
		})
	}
}

func TestClient_DeleteSecret(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		secretName  string
		setupFunc   func(*fake.Clientset)
		wantErr     bool
		errorString string
	}{
		{
			name:       "successful deletion",
			namespace:  "default",
			secretName: "test-secret",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantErr: false,
		},
		{
			name:        "secret not found",
			namespace:   "default",
			secretName:  "non-existing-secret",
			setupFunc:   func(clientset *fake.Clientset) {},
			wantErr:     true,
			errorString: "secret non-existing-secret not found in namespace default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			err := client.DeleteSecret(context.TODO(), tt.namespace, tt.secretName)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("DeleteSecret() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr {
				_, err := client.GetSecret(context.TODO(), tt.namespace, tt.secretName)
				if !errors.IsNotFound(err) {
					t.Errorf("Secret was not deleted properly")
				}
			}
		})
	}
}

func TestClient_GetSecret(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		secretName  string
		setupFunc   func(*fake.Clientset)
		wantErr     bool
		errorString string
	}{
		{
			name:       "secret exists",
			namespace:  "default",
			secretName: "test-secret",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"key1": []byte("value1"),
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantErr: false,
		},
		{
			name:        "secret not found",
			namespace:   "default",
			secretName:  "non-existing-secret",
			setupFunc:   func(clientset *fake.Clientset) {},
			wantErr:     true,
			errorString: "secret non-existing-secret not found in namespace default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			secret, err := client.GetSecret(context.TODO(), tt.namespace, tt.secretName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("GetSecret() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr && secret == nil {
				t.Error("GetSecret() returned nil secret when error was not expected")
			}
		})
	}
}

func TestClient_ListSecrets(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		setupFunc   func(*fake.Clientset)
		wantCount   int
		wantErr     bool
		errorString string
	}{
		{
			name:      "list multiple secrets",
			namespace: "default",
			setupFunc: func(clientset *fake.Clientset) {
				secrets := []corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret1",
							Namespace: "default",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret2",
							Namespace: "default",
						},
					},
				}
				for _, secret := range secrets {
					clientset.CoreV1().Secrets("default").Create(context.TODO(), &secret, metav1.CreateOptions{})
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "empty namespace",
			namespace: "default",
			setupFunc: func(clientset *fake.Clientset) {},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "secrets in different namespaces",
			namespace: "test-ns",
			setupFunc: func(clientset *fake.Clientset) {
				secrets := []corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret1",
							Namespace: "test-ns",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret2",
							Namespace: "different-ns",
						},
					},
				}
				for _, secret := range secrets {
					clientset.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), &secret, metav1.CreateOptions{})
				}
			},
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			secrets, err := client.ListSecrets(context.TODO(), tt.namespace)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("ListSecrets() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr {
				if len(secrets) != tt.wantCount {
					t.Errorf("ListSecrets() returned %d secrets, want %d", len(secrets), tt.wantCount)
				}

				for _, secret := range secrets {
					if secret.Namespace != tt.namespace {
						t.Errorf("ListSecrets() returned secret from wrong namespace: got %s, want %s",
							secret.Namespace, tt.namespace)
					}
				}
			}
		})
	}
}

func TestClient_ListSecretsWithLabels(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		labels      map[string]string
		setupFunc   func(*fake.Clientset)
		wantCount   int
		wantErr     bool
		errorString string
	}{
		{
			name:      "list secrets with matching labels",
			namespace: "default",
			labels: map[string]string{
				"app": "test",
				"env": "prod",
			},
			setupFunc: func(clientset *fake.Clientset) {
				secrets := []corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret1",
							Namespace: "default",
							Labels: map[string]string{
								"app": "test",
								"env": "prod",
							},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret2",
							Namespace: "default",
							Labels: map[string]string{
								"app": "test",
								"env": "dev",
							},
						},
					},
				}
				for _, secret := range secrets {
					clientset.CoreV1().Secrets("default").Create(context.TODO(), &secret, metav1.CreateOptions{})
				}
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "no secrets with matching labels",
			namespace: "default",
			labels: map[string]string{
				"app": "nonexistent",
			},
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret1",
						Namespace: "default",
						Labels: map[string]string{
							"app": "test",
						},
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			selector := metav1.ListOptions{
				LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
					MatchLabels: tt.labels,
				}),
			}

			secrets, err := client.clientset.CoreV1().Secrets(tt.namespace).List(context.TODO(), selector)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("ListSecrets() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr && len(secrets.Items) != tt.wantCount {
				t.Errorf("ListSecrets() returned %d secrets, want %d", len(secrets.Items), tt.wantCount)
			}
		})
	}
}

func TestClient_GetSecretString(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		secretName  string
		key         string
		setupFunc   func(*fake.Clientset)
		want        string
		wantErr     bool
		errorString string
	}{
		{
			name:       "get existing key",
			namespace:  "default",
			secretName: "test-secret",
			key:        "key1",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"key1": []byte("value1"),
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name:       "key not found",
			namespace:  "default",
			secretName: "test-secret",
			key:        "non-existing-key",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"key1": []byte("value1"),
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			wantErr:     true,
			errorString: "key non-existing-key not found in secret test-secret",
		},
		{
			name:        "secret not found",
			namespace:   "default",
			secretName:  "non-existing-secret",
			key:         "key1",
			setupFunc:   func(clientset *fake.Clientset) {},
			wantErr:     true,
			errorString: "secret non-existing-secret not found in namespace default",
		},
		{
			name:       "empty value",
			namespace:  "default",
			secretName: "test-secret",
			key:        "empty-key",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"empty-key": []byte(""),
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			got, err := client.GetSecretString(context.TODO(), tt.namespace, tt.secretName, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecretString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorString {
				t.Errorf("GetSecretString() error = %v, wantString %v", err, tt.errorString)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("GetSecretString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Exists(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		secretName string
		setupFunc  func(*fake.Clientset)
		want       bool
		wantErr    bool
	}{
		{
			name:       "secret exists",
			namespace:  "default",
			secretName: "test-secret",
			setupFunc: func(clientset *fake.Clientset) {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-secret",
						Namespace: "default",
					},
				}
				clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
			},
			want:    true,
			wantErr: false,
		},
		{
			name:       "secret does not exist",
			namespace:  "default",
			secretName: "non-existing-secret",
			setupFunc:  func(clientset *fake.Clientset) {},
			want:       false,
			wantErr:    false,
		},
		{
			name:       "error checking existence",
			namespace:  "",
			secretName: "test-secret",
			setupFunc:  func(clientset *fake.Clientset) {},
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			tt.setupFunc(clientset)

			client := &Client{
				clientset: clientset,
			}

			got, err := client.Exists(context.TODO(), tt.namespace, tt.secretName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
