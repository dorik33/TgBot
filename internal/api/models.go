package api


type CryptoInfo struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	PriceUSD float64 `json:"price_usd"`
	Rank     int     `json:"rank"`
}
