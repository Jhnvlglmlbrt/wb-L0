package models

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestId    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
	PaymentDt    uint32  `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost uint32  `json:"delivery_cost" validate:"required"`
	GoodsTotal   float32 `json:"goods_total" validate:"required"`
	CustomFee    uint32  `json:"custom_fee"`
}

type Item struct {
	ChrtId      uint32  `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       uint16  `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        uint16  `json:"sale" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float32 `json:"total_price" validate:"required"`
	NmId        uint32  `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      uint16  `json:"status" validate:"required"`
}

type Order struct {
	OrderUid          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shard_key" validate:"required"`
	SmId              int      `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

// Validate выполняет проверку структуры Order на соответствие тегам validate.
// Если есть ошибки, возвращается ошибка, иначе nil.
func (o *Order) Validate() error {
	if err := validate.Struct(o); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}
	return nil
}
