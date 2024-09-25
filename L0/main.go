package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"wb_l0/db"
	t "wb_l0/types"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var (
	cache = make(map[string]t.Order)
	mu    sync.RWMutex
)

func kafkaSubscribe(ctx context.Context) {
	var err error
	DB := db.Connect()

	var orders []t.Order
	if tx := DB.Find(&orders); tx.Error != nil {
		log.Fatal("error of get cache: ", err)
	}

	for _, value := range orders {
		cache[value.OrderUID] = value
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "order-service",
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Println("context canceled, stop consumer")
				break
			}
			log.Println("error of reading messege from Kafka: ", err)
			continue
		}

		var order t.Order
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Println("error of unmarshaling json: ", err)
			continue
		}

		err = t.Validate.Struct(order)
		if err != nil {
			log.Println("error of validation: ", err)
			continue
		}

		result := DB.Create(&order)
		if result.Error != nil {
			log.Println("error of writing to db: ", result.Error)
			continue
		}

		mu.Lock()
		cache[order.OrderUID] = order
		mu.Unlock()

		fmt.Println("Order is saved: ", order.OrderUID)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaSubscribe(ctx)

	r := gin.Default()
	r.LoadHTMLFiles("web/templates/model.html")

	r.GET("/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		mu.RLock()
		order, ok := cache[id]
		mu.RUnlock()
		if ok {
			ctx.HTML(http.StatusOK, "model.html", gin.H{"data": order})
		} else {
			ctx.HTML(http.StatusNotFound, "model.html", gin.H{"error": http.StatusNotFound})
		}
	})

	fmt.Println("Server is run...")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server run failed: ", err)
	}
}
