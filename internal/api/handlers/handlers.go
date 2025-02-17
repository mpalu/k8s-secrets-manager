package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
	"github.com/mpalu/k8s-secrets-manager/internal/validator"
)

type Handler struct {
	client k8s.SecretManager
}

func NewHandler(client k8s.SecretManager) *Handler {
	return &Handler{client: client}
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var secretData k8s.SecretData
	if err := json.NewDecoder(r.Body).Decode(&secretData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateSecretData(&secretData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.client.CreateSecret(r.Context(), &secretData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Secret created successfully",
		"name":    secretData.Name,
	})
}

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := vars["namespace"]

	// Check if required parameters are present
	if name == "" || namespace == "" {
		http.Error(w, "name and namespace are required", http.StatusBadRequest)
		return
	}

	secret, err := h.client.GetSecret(r.Context(), namespace, name)
	if err != nil {
		// Handle not found errors separately from internal errors
		if err.Error() == "secret not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(secret)
}

func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	// Check if namespace parameter is present
	if namespace == "" {
		http.Error(w, "namespace is required", http.StatusBadRequest)
		return
	}

	secrets, err := h.client.ListSecrets(r.Context(), namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(secrets)
}

func (h *Handler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var secretData k8s.SecretData
	if err := json.NewDecoder(r.Body).Decode(&secretData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	secretData.Name = name
	if err := h.client.UpdateSecret(r.Context(), &secretData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	namespace := vars["namespace"]

	if err := h.client.DeleteSecret(r.Context(), namespace, name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
