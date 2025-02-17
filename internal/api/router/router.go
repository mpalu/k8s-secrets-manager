package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpalu/k8s-secrets-manager/internal/api/handlers"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
)

func NewRouter(client *k8s.Client) *mux.Router {
	r := mux.NewRouter()
	h := handlers.NewHandler(client)

	// API v1
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Secrets endpoints
	v1.HandleFunc("/secrets", h.CreateSecret).Methods(http.MethodPost)
	v1.HandleFunc("/secrets", h.ListSecrets).Methods(http.MethodGet)
	v1.HandleFunc("/secrets/{name}", h.GetSecret).Methods(http.MethodGet)
	v1.HandleFunc("/secrets/{name}", h.UpdateSecret).Methods(http.MethodPut)
	v1.HandleFunc("/secrets/{name}", h.DeleteSecret).Methods(http.MethodDelete)

	return r
}
