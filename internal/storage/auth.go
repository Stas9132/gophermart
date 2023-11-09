package storage

import (
	"encoding/json"
	"errors"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Auth) UnmarshalJSON(data []byte) error {
	type AuthAlias Auth
	aliasValue := &struct {
		*AuthAlias
	}{
		AuthAlias: (*AuthAlias)(a),
	}
	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	if len(a.Login) < 4 {
		return errors.New("login short")
	}
	if len(a.Password) < 4 {
		return errors.New("password short")
	}
	return nil
}

func (s *DBStorage) RegisterUser(auth Auth) (bool, error) {
	return true, nil
}

func (s *DBStorage) LoginUser(auth Auth) (bool, error) {
	return true, nil
}
