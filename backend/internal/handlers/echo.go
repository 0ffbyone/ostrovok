package handlers

import (
	"backend/pkg/utils"
	"io"
	"net/http"
)

type EchoRequest struct {
	Message string `json:"message"`
}

func (h *handlers) Echo(w http.ResponseWriter, r *http.Request) {
	var size int64 = 0
	if r.ContentLength > 0 {
		size = r.ContentLength
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	clone := make([]byte, 0, size)
	var err error

	defer r.Body.Close()
	clone, err = io.ReadAll(r.Body)

	requestBody := string(clone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, requestBody, http.StatusOK)
}
