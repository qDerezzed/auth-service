// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-service/internal/usecase"
)

func NewRouter(handler *gin.Engine, uc usecase.Auth) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	handler.LoadHTMLGlob("templates/*")
	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, uc)
	}
}
