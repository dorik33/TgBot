package api

type Crypto struct {
	Data `json:"data"`
}

type Data struct {
	ID       string `json:"id"`
	Rank     string `json:"rank"`
	Symbol   string `json:"symbol"`
	PriceUSD string `json:"priceUsd"`
}
