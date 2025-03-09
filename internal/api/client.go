package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type APIClient struct {
	client *http.Client
}

func NewAPIClient(timeout time.Duration) *APIClient {
	return &APIClient{
		&http.Client{
			Timeout: timeout,
		},
	}
}

func (a *APIClient) GetInfo(id string) (Crypto, error) {
	var res Crypto
	url := "https://api.coincap.io/v2/assets/" + id
	resp, err := a.client.Get(url)
	if err != nil {
		log.Println("error while get info from coincap")
		return res, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Println("error while unMarshall info")
		return res, err
	}
	log.Println(res, err)
	if res.PriceUSD == "" {
		return res, fmt.Errorf("Token was not found")
	}

	return res, nil
}
