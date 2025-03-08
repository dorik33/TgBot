package api

import (
	"io"
	"log"
	"net/http"
)

func GetInfo() ([]byte, error) {
	var data []byte
	resp, err := http.Get("https://api.coincap.io/v2/assets")
	if err != nil {
		log.Println("error while get info from coincap")
		return data, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while get info about crypto")
		return data, err
	}
	return data, nil
}


