package api

import (
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"gophermart/internal/auth"
	"gophermart/internal/logger"
	"gophermart/internal/storage"
	"io"
	"net/http"
	"time"
)

type Storage interface {
	io.Closer
	RegisterUser(auth storage.Auth) (bool, error)
	LoginUser(auth storage.Auth) (bool, error)
	NewOrder(order storage.Order) error
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
	w.Write([]byte(auth.GetIssuer(r.Context())))
}

func (h *Handler) PostOrders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := auth.GetIssuer(r.Context())
	if user == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	order, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = goluhn.Validate(string(order)); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = h.storage.NewOrder(storage.Order{
		Number:     string(order),
		Status:     "NEW",
		Accrual:    0,
		UploadedAt: time.Now(),
	})

	if errors.Is(err, storage.ErrSameUser) {
		w.WriteHeader(http.StatusOK)
		return
	} else if errors.Is(err, storage.ErrAnotherUser) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
