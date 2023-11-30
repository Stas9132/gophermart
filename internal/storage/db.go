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
	"time"
)

type DBStorage struct {
	appCtx context.Context
	logger.Logger
	conn *pgx.Conn
	m    map[string]*Order
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
		m:      make(map[string]*Order),
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

var ErrSameUser = errors.New("already in base")
var ErrAnotherUser = errors.New("conflict")

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Issuer     string    `json:"-"`
}

func (s *DBStorage) NewOrder(order Order) error {
	if v, ok := s.m[order.Number]; ok {
		if order.Issuer == v.Issuer {
			return ErrSameUser
		}
		return ErrAnotherUser
	}
	s.m[order.Number] = &order

	//if _, err := s.conn.Exec(s.appCtx, "INSERT INTO orders(number, status, uploaded_at) values ($1,'NEW' ,$2)", order.Number, time.Now()); err != nil {
	//	s.Error("NewOrder() error", slog.String("error", err.Error()))
	//	return err
	//}
	return nil
}

func (s *DBStorage) GetOrders() ([]Order, error) {
	res := make([]Order, 0, len(s.m))
	for _, order := range s.m {
		res = append(res, *order)
	}
	return res, nil
}
