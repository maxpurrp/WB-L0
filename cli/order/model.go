package order

import (
	"encoding/json"
	"fmt"
	"os"
)

type Order struct {
	OrderUID        string `json:"order_uid"`
	TrackNumber     string `json:"track_number"`
	Entry           string `json:"entry"`
	Delivery        DeliveryData
	Payment         PaymentData
	Items           []ItemData
	Locale          string `json:"locale"`
	InternalSig     string `json:"internal_signature"`
	CustomerID      string `json:"customer_id"`
	DeliveryService string `json:"delivery_service"`
	ShardKey        string `json:"shardkey"`
	SMID            int    `json:"sm_id"`
	DateCreated     string `json:"date_created"`
	OofShard        string `json:"oof_shard"`
}

type DeliveryData struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type PaymentData struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float64 `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type ItemData struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	RID         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NMID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

func validation(data []byte) error {
	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}
	return nil
}

func GetOrder() (Order, error) {
	filePath := "order/data/model4.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return Order{}, fmt.Errorf("error reading JSON file: %v", err)
	}

	if err := validation(jsonData); err != nil {
		return Order{}, fmt.Errorf("error validating JSON: %v", err)
	}

	var orderData Order
	err = json.Unmarshal(jsonData, &orderData)
	if err != nil {
		return Order{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return orderData, nil
}
