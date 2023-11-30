package api

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	auth2 "gophermart/internal/auth"
	"gophermart/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

const TTL = time.Hour

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var auth storage.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		h.Error("json decode error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(auth.Login) < 4 || len(auth.Password) < 4 {
		http.Error(w, "invalid login/password", http.StatusBadRequest)
		return
	}

	h.Info("Register request", slog.String("user", auth.Login))

	if ok, err := h.storage.RegisterUser(auth); err != nil {
		h.Error("storage: register user error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		h.Warn("storage: unable create user", slog.String("user", auth.Login))
		http.Error(w, "User already exist", http.StatusConflict)
		return
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": auth.Login, "exp": time.Now().Add(TTL).Unix()},
	).SignedString([]byte(auth2.Key))
	if err != nil {
		h.Error("create jwt error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", j)
	http.SetCookie(w, &http.Cookie{
		Name:   "Authorization",
		Value:  j,
		MaxAge: int(TTL / time.Second),
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var auth storage.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		h.Error("json decode error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Info("Login request", slog.String("user", auth.Login))

	if ok, err := h.storage.LoginUser(auth); err != nil {
		h.Error("storage: logon user error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !ok {
		h.Warn("authentication failed", slog.String("user", auth.Login))
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": auth.Login, "exp": time.Now().Add(TTL).Unix()},
	).SignedString([]byte(auth2.Key))
	if err != nil {
		h.Error("create jwt error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", j)
	http.SetCookie(w, &http.Cookie{
		Name:   "Authorization",
		Value:  j,
		MaxAge: int(TTL / time.Second),
	})

	w.WriteHeader(http.StatusOK)
}
