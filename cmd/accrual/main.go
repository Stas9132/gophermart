package main

import (
	"context"
	"errors"
	"gophermart/internal/accural/api"
	"gophermart/internal/accural/storage"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdownChan := make(chan struct{})
	var server *http.Server

	c := config.New()
	l := logger.NewSlogLogger(c)
	st, err := storage.NewDBStorageAccural(context.Background(), c, l)
	if err != nil {
		log.Fatal(err)
	}

	_ = api.NewAccuralHandler(st, l)

	if err = run(c, server); err != nil {
		panic(err)
	}
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan
		close(shutdownChan)
	}()

	<-shutdownChan
	defer st.Conn.Close(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при завершении работы сервера: %v\n", err)
	}

	os.Exit(0)
}

func run(c *config.Config, server *http.Server) error {
	log.Printf("Сервер запущен на %v\n", c.Address)

	server = &http.Server{Addr: c.Address}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	return nil
}
