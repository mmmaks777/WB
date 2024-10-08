package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"wb_l0/internal/cache"
	"wb_l0/internal/handler"
	"wb_l0/internal/kafka_processor"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error of reading config: %v", err)
	}
	return nil
}

func main() {
	var (
		orderCache *cache.OrderCache
		logger     *zap.Logger
		err        error
	)

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err = initConfig(); err != nil {
		logger.Fatal("error", zap.Error(err))
	}

	orderCache, err = cache.NewOrderCache(1000)
	if err != nil {
		logger.Fatal("error of initialize cache", zap.Error(err))
	}

	shutdownChan := make(chan struct{})
	wg := sync.WaitGroup{}

	wg.Add(1)
	go kafka_processor.KafkaSubscribe(shutdownChan, &wg, orderCache, logger)

	r := gin.Default()
	r.LoadHTMLFiles("internal/templates/model.html")

	handler := handler.OrderHandler{}
	r.GET("/orders/:id", handler.GetOrder(orderCache))

	srv := http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		close(shutdownChan)
		logger.Info("Termination signal received, shutdown the server...")

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Fatal("error when shutting down the server: ", zap.Error(err))
		}
	}()

	logger.Info("Server is run...")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Server run failed: ", zap.Error(err))
	}

	wg.Wait()

	logger.Info("Server stopped successfully")
}
