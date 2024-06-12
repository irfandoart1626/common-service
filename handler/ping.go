package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping Handler provide ping endpoint for OCP/openshift liveness probe
func GetPingHandler() (method, path string, handler gin.HandlerFunc) {
	return "GET", "/_internal/_ping", handlePing
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
