package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/service"
	"gophermart/internal/accural/storage"
	"gophermart/pkg/logger"
	"log"
	"net/http"
)

type StorageAccural interface {
	GetCalculatedDiscountByOrderID(orderID string) (decimal.Decimal, error)
	AcceptOrder(ctx context.Context, order service.Order) error
	AcceptDiscount(ctx context.Context, discount storage.Discount) error
}

type Handler struct {
	logger.Logger
	om StorageAccural
}

func NewAccuralHandler(storage *storage.DBStorage, logger logger.Logger) *Handler {
	return &Handler{
		Logger: logger,
		om:     service.NewOrderManager(storage, logger),
	}

}
func (h Handler) AccrualGoods(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Expected Content-Type: application/json", http.StatusBadRequest)
		return
	}

	var discount storage.Discount
	err := json.NewDecoder(r.Body).Decode(&discount)
	log.Println("+", discount)
	if err != nil {
		h.Error("json.NewDecoder() error", logger.LogMap{"error": err})
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}

	err = h.om.AcceptDiscount(r.Context(), discount)
	if err != nil {
		h.Error("h.om.AcceptDiscount(r.Context(), discount) error", logger.LogMap{"error": err})
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) AccrualOrders(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Expected Content-Type: application/json", http.StatusBadRequest)
		return
	}

	var orderData service.Order
	err := json.NewDecoder(r.Body).Decode(&orderData)
	log.Println("++", orderData)
	if err != nil {
		h.Error("json.NewDecoder() error", logger.LogMap{"error": err})
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}

	err = h.om.AcceptOrder(r.Context(), orderData)
	if err != nil {
		h.Error("h.om.AcceptOrder(r.Context(), orderData) error", logger.LogMap{"error": err})
		http.Error(w, "Order entry error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h Handler) AccrualGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	orderID := vars["number"]

	discount, err := h.om.GetCalculatedDiscountByOrderID(orderID)
	if err != nil {
		h.Error("h.om.GetCalculatedDiscountByOrderID(orderID) error", logger.LogMap{"error": err})
		http.Error(w, "Failed to fetch discount for the order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(struct {
		Order   string          `json:"order"`
		Status  string          `json:"status"`
		Accrual decimal.Decimal `json:"accrual"`
	}{
		Order: orderID,
		Status: func() string {
			if discount.LessThan(decimal.Zero) {
				return "INVALID"
			}
			return "PROCESSED"
		}(),
		Accrual: func() decimal.Decimal {
			if discount.LessThan(decimal.Zero) {
				return decimal.Zero
			}
			return discount
		}(),
	}); err != nil {
		h.Error("json.NewEncoder(w).Encode() error", logger.LogMap{"error": err})
	}
}
