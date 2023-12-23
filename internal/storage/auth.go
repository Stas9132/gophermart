package storage

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"gophermart/pkg/logger"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *DBStorage) RegisterUser(ctx context.Context, auth Auth) (bool, error) {
	h := sha1.Sum([]byte(auth.Password))
	p := hex.EncodeToString(h[:])
	if _, err := s.conn.Exec(ctx, "INSERT INTO auth(login, password) values ($1, $2)", auth.Login, p); err != nil {
		s.Error("unable insert into auth table", logger.LogMap{"error": err, "login": auth.Login})
		return false, err
	}
	return true, nil
}

func (s *DBStorage) LoginUser(ctx context.Context, auth Auth) (bool, error) {
	var p string
	if err := s.conn.QueryRow(ctx, "SELECT password FROM auth where login = $1", auth.Login).Scan(&p); err != nil {
		s.Error("unable select from auth table", logger.LogMap{"error": err, "login": auth.Login})
		return false, err
	}
	ht, err := hex.DecodeString(p)
	if err != nil {
		s.Error("error while decode password", logger.LogMap{"error": err, "login": auth.Login})
		return false, err
	}
	hr := sha1.Sum([]byte(auth.Password))
	if !bytes.Equal(hr[:], ht) {
		s.Warn("Login failure", logger.LogMap{"login": auth.Login})
		return false, err
	}
	return true, nil
}
