package api

import (
	"encoding/json"
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"gophermart/internal/auth"
	"gophermart/internal/logger"
	"gophermart/internal/storage"
	"io"
	"net/http"
	"sort"
	"time"
)

type Storage interface {
	io.Closer
	RegisterUser(auth storage.Auth) (bool, error)
	LoginUser(auth storage.Auth) (bool, error)
	NewOrder(order storage.Order) error
	GetOrders() ([]storage.Order, error)
}

type Handler struct {
	storage Storage
	logger.Logger
}

func NewHandler(storage Storage, l logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		Logger:  l,
	}

}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(auth.GetIssuer(r.Context())))
}

func (h *Handler) PostOrders(w http.ResponseWriter, r *http.Request) {
	var user string
	var order []byte
	defer func() {
		h.Info("POST /api/user/order request", logger.LogMap{"user": user, "order": string(order)})
	}()

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user = auth.GetIssuer(r.Context())
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
		Issuer:     auth.GetIssuer(r.Context()),
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
	var user string
	defer func() {
		h.Info("GET /api/user/orders request", logger.LogMap{"user": user})
	}()
	w.Header().Set("Content-Type", "application/json")
	user = auth.GetIssuer(r.Context())
	if user == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ordrs, err := h.storage.GetOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.SliceStable(ordrs, func(i, j int) bool {
		return ordrs[i].UploadedAt.Before(ordrs[j].UploadedAt)
	})

	if len(ordrs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(ordrs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
