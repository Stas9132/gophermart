package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gophermart/internal/config"
	"gophermart/internal/logger"
)

type DBStorage struct {
	appCtx context.Context
	logger.Logger
}

func NewDBStorage(ctx context.Context, config *config.Config, logger logger.Logger) (*DBStorage, error) {
	createDB(config.DatabaseURI)
	return &DBStorage{
		appCtx: ctx,
		Logger: logger,
	}, nil
}

func (s *DBStorage) Close() error {
	return nil
}

func createDB(DBConn string) {
	conn, err := pgx.Connect(context.Background(), DBConn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database!")

	var tableExists bool

	err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "users").Scan(&tableExists)
	if err != nil {
		panic(err)
	}

	if !tableExists {
		_, err = conn.Exec(context.Background(), `CREATE TABLE users (
	   id SERIAL PRIMARY KEY,
	   name VARCHAR(50),
	   email VARCHAR(50),
	   password VARCHAR(50),
	   points INTEGER
	);`)
		if err != nil {
			panic(err)
		}
		fmt.Println("Table 'users' created.")
	} else {
		fmt.Println("Table 'users' already exist.")
	}

	err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "orders").Scan(&tableExists)
	if err != nil {
		panic(err)
	}

	if !tableExists {
		_, err = conn.Exec(context.Background(), `CREATE TABLE orders (
           order_id SERIAL PRIMARY KEY,
           user_id INTEGER REFERENCES users(id),
           order_date TIMESTAMP,
           total_amount DECIMAL(10, 2),
           status VARCHAR(50)
        );`)
		if err != nil {
			panic(err)
		}
		fmt.Println("Table 'orders' created.")
	} else {
		fmt.Println("Table 'orders' already exists.")
	}
}
