package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
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
			return []byte(Key), nil
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
		next.ServeHTTP(w, r.WithContext(context.WithoutCancel(ctx)))
	})
}

const Key = "secret_key"
