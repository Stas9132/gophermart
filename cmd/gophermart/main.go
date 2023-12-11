package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"gophermart/internal/api"
	"gophermart/internal/auth"
	"gophermart/internal/config"
	"gophermart/internal/logger"
	"gophermart/internal/storage"
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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c := config.New()
	l := logger.NewSlogLogger(c)
	st, err := storage.NewDBStorage(ctx, c, l)
	if err != nil {
		log.Fatal("storage open error", err)
	}
	h := api.NewHandler(st, l)

	mRouter(h)
	server = &http.Server{Addr: c.Address}
	run(c)

	<-ctx.Done()

	if err = server.Close(); err != nil {
		log.Println("server close error", err)
	}

	time.Sleep(time.Second)
	os.Exit(0)
}

func mRouter(handler *api.Handler) {
	r := mux.NewRouter()

	r.Use(auth.Authorization)

	r.HandleFunc("/api/user/test", handler.Test).Methods(http.MethodGet)
	r.HandleFunc("/api/user/register", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", handler.Login).Methods(http.MethodPost)

	r.HandleFunc("/api/user/orders", handler.PostOrders).Methods(http.MethodPost)
	r.HandleFunc("/api/user/orders", handler.GetOrders).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance", handler.GetBalance).Methods(http.MethodGet)
	r.HandleFunc("/api/user/balance/withdraw", handler.PostBalanceWithdraw).Methods(http.MethodPost)
	r.HandleFunc("/api/user/withdrawals", handler.GetWithdraw).Methods(http.MethodGet)

	http.Handle("/", r)
}
