package api

import "net/http"

func (h *Handler) AccrualGoods(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AccrualOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) AccrualGetOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
