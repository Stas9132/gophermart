package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v5"
	"gophermart/internal/config"
	"gophermart/internal/logger"
	"log/slog"
)

type DBStorage struct {
	appCtx context.Context
	db     *sql.DB
	logger.Logger
}

func NewDBStorage(ctx context.Context, config *config.Config, logger logger.Logger) (*DBStorage, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=user1 password=user1 dbname=user1 sslmode=disable")
	if err != nil {
		logger.Error("open database failed", slog.String("error", err.Error()))
		return nil, err
	}
	return &DBStorage{
		appCtx: ctx,
		db:     db,
		Logger: logger,
	}, nil
}

func (s *DBStorage) Close() error {
	return nil
}
