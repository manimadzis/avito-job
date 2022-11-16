package main

import (
	"avito-job/internal/app"
	"avito-job/internal/config"
	"avito-job/pkg/logging"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	logger := logging.Get()
	conf := &config.Config{
		ServerHost:   "localhost",
		ServerPort:   "9876",
		DBHost:       "localhost",
		DBPort:       "5432",
		DBUsername:   "postgres",
		DBPassword:   "pass",
		DatabaseName: "bank_service",
	}

	a := app.NewApp(conf, *logger)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt)
		<-s
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if err := a.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error occured while shuting down server: %v", err)
		} else if err == context.Canceled {
			cancel()
			logger.Fatalf("Can't close server. Timeout")
		}

		logger.Info("Shutdown server")
		wg.Done()
	}()

	go func() {
		err := a.Start()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("%v", err)
		}
		wg.Done()
	}()

	wg.Wait()
}
