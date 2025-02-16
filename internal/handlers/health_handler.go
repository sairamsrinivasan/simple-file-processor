package handlers

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler handles the health check request
func (h handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
