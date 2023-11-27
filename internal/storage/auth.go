package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"log/slog"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *DBStorage) RegisterUser(auth Auth) (bool, error) {
	h := sha1.Sum([]byte(auth.Password))
	p := hex.EncodeToString(h[:])
	if _, err := s.conn.Exec(s.appCtx, "INSERT INTO auth(login, password) values ($1, $2)", auth.Login, p); err != nil {
		s.Error("unable insert into auth table", slog.String("error", err.Error()), slog.String("login", auth.Login))
		return false, err
	}
	return true, nil
}

func (s *DBStorage) LoginUser(auth Auth) (bool, error) {
	var p string
	if err := s.conn.QueryRow(s.appCtx, "SELECT password FROM auth where login = $1", auth.Login).Scan(&p); err != nil {
		s.Error("unable select from auth table", slog.String("error", err.Error()), slog.String("login", auth.Login))
		return false, err
	}
	ht, err := hex.DecodeString(p)
	if err != nil {
		s.Error("error while decode password", slog.String("error", err.Error()), slog.String("login", auth.Login))
		return false, err
	}
	hr := sha1.Sum([]byte(auth.Password))
	if !bytes.Equal(hr[:], ht) {
		s.Warn("Login failure", slog.String("login", auth.Login))
		return false, err
	}
	return true, nil
}
