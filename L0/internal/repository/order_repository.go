package repository

import (
	"fmt"
	"wb_l0/internal/cache"
	"wb_l0/internal/domain"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	user := viper.GetString("database.user")
	dbname := viper.GetString("database.name")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	dsn := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=disable", user, dbname, host, port)
	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error of connection to DB: %v", err)
	}

	err = DB.AutoMigrate(&domain.Order{})
	if err != nil {
		return nil, fmt.Errorf("error of auto migrate: %v", err)
	}

	return DB, nil
}

func UploadCache(db *gorm.DB, orderCache *cache.OrderCache) error {
	var orders []domain.Order
	if tx := db.Find(&orders); tx.Error != nil {
		return fmt.Errorf("error of get cache: %v", tx.Error)
	}

	for _, value := range orders {
		orderCache.Set(value)
	}

	return nil
}
