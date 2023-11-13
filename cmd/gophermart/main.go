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
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	log.Println("Server starting")
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
		server = &http.Server{Addr: c.Host + ":" + c.Port}
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

	r.HandleFunc("/test", handler.Test).Methods("GET")

	http.Handle("/", r)
}
