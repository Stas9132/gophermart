package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/shopspring/decimal"
	"gophermart/internal/auth"
	"gophermart/internal/logger"
	"gophermart/internal/storage"
	"io"
	"net/http"
	"sort"
	"time"
)

type Storage interface {
	RegisterUser(ctx context.Context, auth storage.Auth) (bool, error)
	LoginUser(ctx context.Context, auth storage.Auth) (bool, error)
	NewOrder(ctx context.Context, order storage.Order) error
	GetOrders(ctx context.Context) ([]storage.Order, error)
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
	_, _ = w.Write([]byte(auth.GetIssuer(r.Context())))
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
	err = h.storage.NewOrder(r.Context(), storage.Order{
		Number:     string(order),
		Status:     "NEW",
		Accrual:    decimal.Zero,
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

	ordrs, err := h.storage.GetOrders(r.Context())
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

var balance = struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}{
	Current:   decimal.NewFromFloat32(729.98),
	Withdrawn: decimal.Zero,
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(balance)
}

func (h *Handler) PostBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Order string `json:"order"`
		Sum   int    `json:"sum"`
	}
	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Error("json.Decode()", logger.LogMap{"error": err})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = goluhn.Validate(req.Order); err != nil {
		h.Error("goluhn.Validate()", logger.LogMap{"error": err, "order": req.Order})
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	balance.Current = balance.Current.Sub(decimal.NewFromInt(int64(req.Sum)))
	balance.Withdrawn = balance.Withdrawn.Add(decimal.NewFromInt(int64(req.Sum)))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
