package api

import (
	"gophermart/internal/logger"
	"gophermart/internal/storage"
	"io"
	"net/http"
)

type Storage interface {
	io.Closer
	RegisterUser(auth storage.Auth) (bool, error)
	LoginUser(auth storage.Auth) (bool, error)
}

type Handler struct {
	storage Storage
	logger.Logger
}

func NewHandler(storage Storage, logger logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		Logger:  logger,
	}

}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

}
