package main

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"gophermart/internal/api"
	"gophermart/internal/app/process"
	"gophermart/internal/storage"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	server *http.Server
)

func run(c *config.Config) {
	log.Println("Server starting")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			} else {
				log.Fatal(err)
			}
		}
	}()
}

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c := config.New()
	l := logger.NewSlogLogger(c)
	st, err := storage.NewDBStorage(ctx, c, l)
	if err != nil {
		log.Fatal("storage open error", err)
	}
	api.NewHandler(st, l)
	go process.StatusDaemon(ctx, c, st, l)
	server = &http.Server{Addr: c.Address}
	run(c)

	<-ctx.Done()

	if err = server.Close(); err != nil {
		log.Println("server close error", err)
	}

	time.Sleep(time.Second)
	os.Exit(0)
}
