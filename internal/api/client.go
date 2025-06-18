package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dorik33/TgBot/internal/models"
)

type aPIClient struct {
	APIKey string
	Client *http.Client
}

type CoinAPI interface {
	GetInfo(symbol string) (*models.CryptoInfo, error)
}

func NewAPIClient(apiKey string, timeout time.Duration) CoinAPI {
	return &aPIClient{
		APIKey: apiKey,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *aPIClient) GetInfo(symbol string) (*models.CryptoInfo, error) {
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s&convert=USD", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accepts", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result struct {
		Data map[string]struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
			CmcRank int    `json:"cmc_rank"`
			Quote   map[string]struct {
				Price float64 `json:"price"`
			} `json:"quote"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	data := result.Data[strings.ToUpper(symbol)]
	info := &models.CryptoInfo{
		Name:     data.Name,
		Symbol:   data.Symbol,
		PriceUSD: data.Quote["USD"].Price,
	}

	return info, nil
}
