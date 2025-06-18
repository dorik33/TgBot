package repository

import "github.com/dorik33/TgBot/internal/models"

type SubscriptionRepository interface {
	AddSubscription(int64, string) error
	GetSubcriptions(int64) ([]models.Subscription, error)
	DeleteSubscription(int64, string) error
	GetAllSubs() ([]models.Subscription, error)
}

type WalletRepository interface {
	GetWallet(int64) ([]models.Portfolio, error)
	AddCrypto(*models.Portfolio) error
	DeleteCrypto(int, int64) error
}

type RedisRepository interface {
	SetCryptoPrice(string, string) error
	GetCryptoPrice(string) (string, error)
	Close()
}