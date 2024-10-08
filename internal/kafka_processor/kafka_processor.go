package kafka_processor

import (
	"context"
	"encoding/json"
	"fmt"
	"wb_l0/internal/cache"
	"wb_l0/internal/domain"
	"wb_l0/internal/repository"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KafkaMessageProcessor interface {
	ReadMessage(ctx context.Context) (domain.Order, error)
	Validate(order domain.Order) error
	UpdateCache(order domain.Order) error
	SaveToDB(order domain.Order) error
}

type KafkaProcessor struct {
	Reader *kafka.Reader
	DB     *gorm.DB
	logger *zap.Logger
}

func NewKafkaProcessor(logger *zap.Logger) KafkaProcessor {
	kp := KafkaProcessor{}

	kp.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: viper.GetStringSlice("kafka.brokers"),
		Topic:   viper.GetString("kafka.topic"),
		GroupID: viper.GetString("kafka.groupid"),
	})

	var err error
	kp.DB, err = repository.Connect()
	if err != nil {
		logger.Fatal("error", zap.Error(err))
	}

	kp.logger = logger

	return kp
}

func (p *KafkaProcessor) ReadMessage() (domain.Order, error) {
	msg, err := p.Reader.ReadMessage(context.Background())
	if err != nil {
		return domain.Order{}, err
	}

	var order domain.Order
	if err = json.Unmarshal(msg.Value, &order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (p *KafkaProcessor) Validate(order domain.Order) error {
	return domain.Validate.Struct(order)
}

func (p *KafkaProcessor) UpdateCache(order domain.Order, cache *cache.OrderCache) error {
	if _, exists := cache.Get(order.OrderUID); exists {
		return fmt.Errorf("dublicate key: %s", order.OrderUID)
	}

	cache.Set(order)
	return nil
}

func (p *KafkaProcessor) SaveToDB(order domain.Order) error {
	return p.DB.Create(&order).Error
}
