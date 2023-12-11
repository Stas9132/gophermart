package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"gophermart/internal/accural/api"
	"gophermart/internal/config"
	"gophermart/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	server       *http.Server
	shutdownChan = make(chan struct{})
)

func main() {

	c := config.New()
	l := logger.NewSlogLogger(c)
	//st, err := storage.NewDBStorageAccural(context.Background(), c, l)
	//if err != nil {
	//	log.Println(err)
	//
	//}
	handler := api.NewAccuralHandler(nil, l)

	mRouter(handler)
	if err := run(c); err != nil {
		panic(err)
	}
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan

		close(shutdownChan)
	}()

	<-shutdownChan
	//defer st.Conn.Close(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении работы сервера: %v\n", err)
	}

	os.Exit(0)
}

func run(c *config.Config) error {
	log.Printf("Сервер запущен на %v\n", c.Address)

	server = &http.Server{Addr: c.Address}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-shutdownChan
	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении работы сервера: %v\n", err)
	}

	return nil
}

func mRouter(handler *api.Handler) {
	r := mux.NewRouter()

	r.HandleFunc("/api/goods", handler.AccrualGoods).Methods(http.MethodPost)
	r.HandleFunc("/api/orders", handler.AccrualOrders).Methods(http.MethodPost)
	r.HandleFunc("/api/orders/{number}", handler.AccrualGetOrders).Methods(http.MethodGet)
	http.Handle("/", r)
}
