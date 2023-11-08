package storage

import (
	"context"
	"gophermart/internal/config"
	"gophermart/internal/logger"
)

type DBStorage struct {
	appCtx context.Context
	logger.Logger
}

func NewDBStorage(ctx context.Context, config *config.Config, logger logger.Logger) (*DBStorage, error) {
	return &DBStorage{
		appCtx: ctx,
		Logger: logger,
	}, nil
}

func (s *DBStorage) Close() error {
	return nil
}
