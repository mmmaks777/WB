package main

import (
	"context"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("error of reading config", zap.Error(err))
	}
}

func main() {
	initConfig()
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: viper.GetStringSlice("kafka.brokers"),
		Topic:   viper.GetString("kafka.topic"),
	})
	defer writer.Close()

	messege, err := os.ReadFile("testdata/model.json")
	if err != nil {
		logger.Fatal("error of reading json file: ", zap.Error(err))
	}

	msg := kafka.Message{
		Key:   []byte("order-" + time.Now().Format("20060102150405")),
		Value: messege,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		logger.Fatal("error writing messege to kafka: ", zap.Error(err))
	}

	logger.Info("Messege published")
}
