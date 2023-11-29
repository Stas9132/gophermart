package storage

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gophermart/internal/config"
	"gophermart/internal/logger"
	"log/slog"
)

type DBStorage struct {
	appCtx context.Context
	logger.Logger
	conn *pgx.Conn
}

func NewDBStorage(ctx context.Context, config *config.Config, logger logger.Logger) (*DBStorage, error) {
	conn, err := createDB(config.DatabaseURI, logger)
	if err != nil {
		return nil, err
	}
	return &DBStorage{
		appCtx: ctx,
		Logger: logger,
		conn:   conn,
	}, nil
}

func (s *DBStorage) Close() error {
	return nil
}

func createDB(DBConn string, logger logger.Logger) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), DBConn)
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to the database!", slog.String("DSN", DBConn))

	m, err := migrate.New("file://internal/storage/migration", DBConn)
	if err != nil {
		logger.Error("Error while create migration", slog.String("error", err.Error()))
		return nil, err
	}
	m.Drop()
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Error while migration up", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Info("Migration complete!")

	return conn, nil
}
