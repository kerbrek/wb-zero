package model

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

type Order struct {
	Uid               string   `json:"order_uid" validate:"required,alphanum,max=50"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required,min=1"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	Shardkey          string   `json:"shardkey" validate:"required"`
	SmId              uint64   `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       uint64 `json:"amount" validate:"required"`
	PaymentDt    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost uint64 `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   uint64 `json:"goods_total" validate:"required"`
	CustomFee    uint64 `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtId      uint64 `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       uint64 `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        uint64 `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  uint64 `json:"total_price" validate:"required"`
	NmId        uint64 `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
