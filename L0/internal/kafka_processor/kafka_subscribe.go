package kafka_processor

import (
	"errors"
	"io"
	"sync"
	"wb_l0/internal/cache"
	"wb_l0/internal/repository"

	"go.uber.org/zap"
)

func KafkaSubscribe(shutdownChan chan struct{}, wg *sync.WaitGroup, orderCache *cache.OrderCache, logger *zap.Logger) {
	defer wg.Done()

	kp := NewKafkaProcessor(logger)

	if err := repository.UploadCache(kp.DB, orderCache); err != nil {
		logger.Error("error", zap.Error(err))
	}

	defer kp.Reader.Close()

	go func() {
		<-shutdownChan
		logger.Info("Shutdown signal received, stopping Kafka consumer...")
		kp.Reader.Close()
	}()

	for {
		order, err := kp.ReadMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Error("error of reading messege from Kafka: ", zap.Error(err))
			continue
		}

		if err := kp.Validate(order); err != nil {
			logger.Error("error of validation: ", zap.Error(err))
			continue
		}

		if err := kp.UpdateCache(order, orderCache); err != nil {
			logger.Error("error: duplicate key: ", zap.String("orderUID", order.OrderUID))
			continue
		}

		if err := kp.SaveToDB(order); err != nil {
			logger.Error("error of writing to db: ", zap.Error(err))
		}

		logger.Info("Order is saved: ", zap.String("orderUID", order.OrderUID))
	}
}
