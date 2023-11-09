package api

import (
	"encoding/json"
	"gophermart/internal/storage"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var auth storage.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok, err := h.storage.RegisterUser(auth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var auth storage.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ok, err := h.storage.LoginUser(auth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
