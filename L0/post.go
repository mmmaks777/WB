package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})
	defer writer.Close()

	messege, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal("error of reading json file: ", err)
	}

	msg := kafka.Message{
		Key:   []byte("order-" + time.Now().Format("20060102150405")),
		Value: messege,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Fatal("error writing messege to kafka: ", err)
	}

	fmt.Println("Messege is publishing")
}
