package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"gophermart/internal/api"
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
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		} else {
			log.Fatal(err)
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c := config.New()
	l := logger.NewSlogLogger(c)
	var err error
	var st api.Storage
	if st, err = storage.NewDBStorage(ctx, c, l); err != nil {
		log.Fatal("storage open error", err)
	}
	h := api.NewHandler(st, l)

	go func() {
		mRouter(h)
		server = &http.Server{Addr: c.Address}
		run(c)
	}()

	<-ctx.Done()

	if err = server.Close(); err != nil {
		log.Println("server close error", err)
	}
	if err = st.Close(); err != nil {
		log.Println("storage close error", err)
	}

	time.Sleep(time.Second)
	os.Exit(0)
}

func mRouter(handler *api.Handler) {
	r := mux.NewRouter()

	//r.Use(handler.LoggingMiddleware, gzip.GzipMiddleware, handler.HashSHA256Middleware)
	//r.Use(api.Authorization)

	r.HandleFunc("/api/goods", handler.AccrualGoods).Methods(http.MethodPost)
	r.HandleFunc("/api/orders", handler.AccrualOrders).Methods(http.MethodPost)

	http.Handle("/", r)
}
