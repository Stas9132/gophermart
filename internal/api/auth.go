package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gophermart/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type issuer struct {
}

func GetIssuer(ctx context.Context) string {
	s, ok := ctx.Value(issuer{}).(string)
	if !ok {
		return ""
	}
	return s
}

const key = "secret_key"
const TTL = time.Hour

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var auth storage.Auth
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		h.Error("json decode error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	).SignedString([]byte(key))
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
	).SignedString([]byte(key))
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

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		j := r.Header.Get("Authorization")
		if c, err := r.Cookie("Authorization"); err == nil {
			j = c.Value
		}

		t, err := jwt.ParseWithClaims(j, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signed method")
			}
			return []byte(key), nil
		})
		if err != nil {
			goto lNext
		}

		if c, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
			if u, ok := (*c)["iss"].(string); ok {
				ctx = context.WithValue(ctx, issuer{}, u)
			}
		}
	lNext:
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
