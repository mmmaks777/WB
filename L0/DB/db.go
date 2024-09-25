package db

import (
	"log"

	"wb_l0/types"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "user=Maksimka dbname=wb_l0 port=5432 sslmode=disable"
	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connection to DB", err)
	}

	err = DB.AutoMigrate(&types.Order{})
	if err != nil {
		log.Fatal("error auto migrate", err)
	}

	return DB
}
