package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/models"
	"gophermart/internal/accural/storage"
	"gophermart/internal/logger"
	"io"
	"net/http"
)

type StorageAccural interface {
	io.Closer
	AcceptOrder(discounts []storage.Discount) error
	AcceptDiscount(discounts []storage.Discount) error
	CalculateDiscount(ds []storage.Discount) (decimal.Decimal, error)
	GetCalculatedDiscountByOrderID(orderID int) (decimal.Decimal, error)
}

type Handler struct {
	storage StorageAccural
	logger.Logger
	om *models.OrderManager
}

func NewAccuralHandler(storage *storage.DBStorage, logger logger.Logger) *Handler {
	return &Handler{
		storage: storage,
		Logger:  logger,
		om:      models.NewOrderManager(storage),
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
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}
	err = models.OrderManager{}.AcceptDiscount(discount)
	if err != nil {
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

	var orderData models.Order
	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}

	err = models.OrderManager{}.AcceptOrder(orderData)
	if err != nil {
		http.Error(w, "Order entry error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h Handler) AccrualGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	vars := mux.Vars(r)
	orderID := vars["number"]

	discount, err := models.OrderManager{}.GetCalculatedDiscountByOrderID(orderID)
	if err != nil {
		http.Error(w, "Failed to fetch discount for the order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(discount.String()))
}
