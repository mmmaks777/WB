package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

type Delivery struct {
	Name    string `validate:"required"`
	Phone   string `validate:"required"`
	Zip     string `validate:"required"`
	City    string `validate:"required"`
	Address string `validate:"required"`
	Region  string `validate:"required"`
	Email   string `validate:"email"`
}

type Payment struct {
	Transaction  string `validate:"required"`
	RequestID    string
	Currency     string `validate:"required"`
	Provider     string `validate:"required"`
	Amount       int    `validate:"required"`
	PaymentDT    int64  `json:"payment_dt" validate:"required"`
	Bank         string `validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required"`
	GoodsTotal   int    `json:"goods_total" validate:"required"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required"`
	NmID        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

type Order struct {
	gorm.Model
	OrderUID          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `validate:"required"`
	Delivery          Delivery `gorm:"embedded;embeddedPrefix:delivery_" validate:"required"`
	Payment           Payment  `gorm:"embedded;embeddedPrefix:payment_" validate:"required"`
	Items             Items    `json:"items" gorm:"type:json" validate:"dive"`
	Locale            string   `validate:"required"`
	InternalSignature string
	CustomerID        string `json:"customer_id"  validate:"required"`
	DeliveryService   string `json:"delivery_service" validate:"required"`
	ShardKey          string `validate:"required"`
	SmID              int    `json:"sm_id" validate:"required"`
	DateCreated       string `json:"date_created" validate:"required"`
	OofShard          string `json:"oof_shard" validate:"required"`
}

type Items []Item

func (items *Items) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("не удалось преобразовать значение %v в []byte", value)
	}
	return json.Unmarshal(bytes, items)
}

func (items Items) Value() (driver.Value, error) {
	return json.Marshal(items)
}
