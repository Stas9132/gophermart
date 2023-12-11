package storage

import (
	"context"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"gophermart/internal/config"
	"gophermart/internal/logger"
)

type DBStorage struct {
	appCtx context.Context
	logger.Logger
	Conn *pgx.Conn
}

func (D DBStorage) Close() error {
	//TODO implement me
	panic("implement me")
}

func (D DBStorage) AcceptOrder(discounts []Discount) error {
	//TODO implement me
	panic("implement me")
}

func (D DBStorage) AcceptDiscount(discounts []Discount) error {
	//TODO implement me
	panic("implement me")
}

func (D DBStorage) CalculateDiscount(ds []Discount) (decimal.Decimal, error) {
	//TODO implement me
	panic("implement me")
}

func (D DBStorage) GetCalculatedDiscountByOrderID(orderID int) (decimal.Decimal, error) {
	//TODO implement me
	panic("implement me")
}

func NewDBStorageAccural(ctx context.Context, config *config.Config, logger logger.Logger) (*DBStorage, error) {
	conn, err := newMigrate(config.DatabaseURI, logger)
	if err != nil {
		return nil, err
	}
	return &DBStorage{
		appCtx: ctx,
		Logger: logger,
		Conn:   conn,
	}, nil
}

func newMigrate(DBConn string, l logger.Logger) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), DBConn)
	if err != nil {
		return nil, err
	}

	l.Info("Successfully connected to the database!", logger.LogMap{"DSN": DBConn})
	return conn, nil
}

type Discount struct {
	Match      string
	Reward     decimal.Decimal
	RewardType string
}
