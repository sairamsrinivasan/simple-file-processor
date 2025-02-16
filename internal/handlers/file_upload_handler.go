package handlers

import (
	"fmt"
	"net/http"
)

// FileUploadHandler handles the file upload request
func (h handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File upload request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
