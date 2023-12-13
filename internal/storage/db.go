package storage

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"gophermart/internal/auth"
	"gophermart/internal/config"
	l2 "gophermart/internal/logger"
	"log"
	"time"
)

type DBStorage struct {
	l2.Logger
	conn *pgx.Conn
	m    map[string]*Order
}

type StorageImpl interface {
	NewOrder(ctx context.Context, order Order) error
	GetOrders(ctx context.Context) ([]Order, error)
	UpdateOrder(ctx context.Context, order Order) error
	GetOrdersInProcessing() ([]Order, error)
}

func New() StorageImpl {
	return &DBStorage{}
}

func NewDBStorage(ctx context.Context, config *config.Config, logger l2.Logger) (*DBStorage, error) {
	conn, err := createDB(config.DatabaseURI, logger)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, "select number, status, accrual, uploaded_at, issuer from orders")
	if err != nil {
		logger.Error("select request error", l2.LogMap{"error": err})
		return nil, err
	}
	m := make(map[string]*Order)
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.Issuer)
		m[order.Number] = &order
		if err != nil {
			logger.Error("scan error", l2.LogMap{"error": err})
			return nil, err
		}
	}
	return &DBStorage{
		Logger: logger,
		conn:   conn,
		m:      m,
	}, nil
}

func (s *DBStorage) Close() error {
	_ = s.conn.Close(context.Background())
	return nil
}

func createDB(DBConn string, logger l2.Logger) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), DBConn)
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to the database!", l2.LogMap{"DSN": DBConn})

	m, err := migrate.New("file://internal/storage/migration", DBConn)
	if err != nil {
		logger.Error("Error while create migration", l2.LogMap{"error": err})
		return nil, err
	}
	_ = m.Drop()
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Error while migration up", l2.LogMap{"error": err})
		return nil, err
	}
	logger.Info("Migration complete!")

	return conn, nil
}

var ErrSameUser = errors.New("already in base")
var ErrAnotherUser = errors.New("conflict")

type Order struct {
	Number     string          `json:"number"`
	Status     string          `json:"status"`
	Accrual    decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt time.Time       `json:"uploaded_at"`
	Issuer     string          `json:"-"`
}

func (s *DBStorage) NewOrder(ctx context.Context, order Order) error {

	if v, ok := s.m[order.Number]; ok {
		if order.Issuer == v.Issuer {
			return ErrSameUser
		}
		return ErrAnotherUser
	}

	s.m[order.Number] = &order

	if _, err := s.conn.Exec(ctx, "INSERT INTO orders (number, status, accrual, uploaded_at, issuer) VALUES ($1,$2,$3,$4,$5);", order.Number, order.Status, order.Accrual, order.UploadedAt, order.Issuer); err != nil {
		s.Error("NewOrder() error", l2.LogMap{"error": err})
		return err
	}

	return nil
}

func (s *DBStorage) GetOrders(ctx context.Context) ([]Order, error) {
	res := make([]Order, 0, len(s.m))
	issuer := auth.GetIssuer(ctx)
	for _, order := range s.m {
		if order.Issuer == issuer {
			res = append(res, *order)
		}
	}
	return res, nil
}
func (s *DBStorage) GetOrdersInProcessing() ([]Order, error) {
	log.Println("cash size :", len(s.m))
	res := make([]Order, 0, len(s.m))
	for _, order := range s.m {
		if order.Status == "NEW" {
			res = append(res, *order)
		}
	}
	log.Println(res)

	return res, nil
}

func (s *DBStorage) UpdateOrder(ctx context.Context, order Order) error {
	_, ok := s.m[order.Number]
	if !ok {
		return errors.New("order not found")
	}

	_, err := s.conn.Exec(context.Background(), "UPDATE orders SET status = $1, accrual = $2 WHERE number = $3;", order.Status, order.Accrual, order.Number)
	if err != nil {
		s.Error("Update Order error", l2.LogMap{"error": err})
		return err
	}

	return nil
}
