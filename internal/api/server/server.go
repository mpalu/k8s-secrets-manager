package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mpalu/k8s-secrets-manager/internal/k8s"
)

type Server struct {
	router *gin.Engine
	client *k8s.Client
}

func New(client *k8s.Client) *Server {
	router := gin.Default()
	server := &Server{
		router: router,
		client: client,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/secrets", s.createSecret)
	s.router.PUT("/secrets/:namespace/:name", s.updateSecret)
	s.router.DELETE("/secrets/:namespace/:name", s.deleteSecret)
	s.router.GET("/secrets/:namespace", s.listSecrets)
	s.router.GET("/secrets/:namespace/:name", s.getSecret)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) createSecret(c *gin.Context) {
	var secretData k8s.SecretData
	if err := c.ShouldBindJSON(&secretData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.client.CreateSecret(c.Request.Context(), &secretData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Secret %s created successfully in namespace %s",
			secretData.Name, secretData.Namespace),
	})
}

func (s *Server) updateSecret(c *gin.Context) {
	var secretData k8s.SecretData
	if err := c.ShouldBindJSON(&secretData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secretData.Name = c.Param("name")
	secretData.Namespace = c.Param("namespace")

	if err := s.client.UpdateSecret(c.Request.Context(), &secretData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Secret %s updated successfully in namespace %s",
			secretData.Name, secretData.Namespace),
	})
}

func (s *Server) deleteSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if err := s.client.DeleteSecret(c.Request.Context(), namespace, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Secret %s deleted successfully from namespace %s",
			name, namespace),
	})
}

func (s *Server) listSecrets(c *gin.Context) {
	namespace := c.Param("namespace")

	secrets, err := s.client.ListSecrets(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, secrets)
}

func (s *Server) getSecret(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	secret, err := s.client.GetSecret(c.Request.Context(), namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, secret)
}
