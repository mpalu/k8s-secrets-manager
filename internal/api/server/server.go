package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mpalu/k8s-secrets-manager/internal/api/handlers"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
)

type Server struct {
	router *mux.Router
	client *k8s.Client
}

func New(client *k8s.Client) *Server {
	router := mux.NewRouter()
	h := handlers.NewHandler(client)

	// API v1
	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/secrets", h.CreateSecret).Methods(http.MethodPost)
	v1.HandleFunc("/secrets", h.ListSecrets).Methods(http.MethodGet)
	v1.HandleFunc("/secrets/{namespace}/{name}", h.GetSecret).Methods(http.MethodGet)
	v1.HandleFunc("/secrets/{namespace}/{name}", h.UpdateSecret).Methods(http.MethodPut)
	v1.HandleFunc("/secrets/{namespace}/{name}", h.DeleteSecret).Methods(http.MethodDelete)

	return &Server{
		router: router,
		client: client,
	}
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
