package handlers

import (
	"backend/pkg/logger"
	"net/http"
)

type Handlers interface {
	Echo(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(logger logger.Logger) Handlers {
	return &handlers{
		logger: logger,
	}
}

type handlers struct {
	logger logger.Logger
}
