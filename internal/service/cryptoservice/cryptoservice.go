package cryptoservice

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/models"
	"github.com/dorik33/TgBot/internal/repository"
	"github.com/dorik33/TgBot/internal/service"
	"github.com/redis/go-redis/v9"
)

type cryptoService struct {
	apiClient api.CoinAPI
	redis     repository.RedisRepository
}

func NewCryptoService(client api.CoinAPI, r repository.RedisRepository) service.CryptoService {
	return &cryptoService{
		apiClient: client,
		redis:     r,
	}
}

func (s *cryptoService) GetCryptoPrice(symbol string) (*models.CryptoInfo, error) {

	cachedData, err := s.redis.GetCryptoPrice(symbol)
	
	if err == nil {
		cryptoInfo, err := deserializeJSON(cachedData)
		if err == nil {
			fmt.Println("Данные получены из Redis")
			return &cryptoInfo, nil
		} else {
			log.Printf("Ошибка при десериализации из Redis: %v", err)
		}
	} else if err != redis.Nil {
		log.Printf("Ошибка при получении данных из Redis: %v", err)
	}

	apiInfo, err := s.apiClient.GetInfo(symbol)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных из API: %w", err)
	}

	go func() {
		jsonData, err := serializeJSON(*apiInfo)
		if err != nil {
			log.Printf("Ошибка при сериализации в JSON для Redis: %v", err)
			return
		}
		if err := s.redis.SetCryptoPrice(symbol, jsonData); err != nil {
			log.Printf("Ошибка при сохранении в Redis: %v", err)
		} else {
			fmt.Println("Данные сохранены в Redis")
		}
	}()

	return apiInfo, nil
}

func serializeJSON(cryptoInfo models.CryptoInfo) (string, error) {
	jsonData, err := json.Marshal(cryptoInfo)
	if err != nil {
		return "", fmt.Errorf("ошибка при сериализации в JSON: %w", err)
	}
	return string(jsonData), nil
}

func deserializeJSON(jsonData string) (models.CryptoInfo, error) {
	var cryptoInfo models.CryptoInfo
	err := json.Unmarshal([]byte(jsonData), &cryptoInfo)
	if err != nil {
		return models.CryptoInfo{}, fmt.Errorf("ошибка при десериализации из JSON: %w", err)
	}
	return cryptoInfo, nil
}
