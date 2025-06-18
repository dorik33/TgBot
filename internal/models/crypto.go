package models

import "time"

type CryptoInfo struct {
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	PriceUSD float64 `json:"price_usd"`
}

type Subscription struct {
	ID         int
	UserID     int64
	Token      string
	Created_at time.Time
}

type Portfolio struct {
	ID         int
	UserID     int64
	Token      string
	Amount     float64
	Price      float64
	Created_at time.Time
}
