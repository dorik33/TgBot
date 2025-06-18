package service

import "github.com/dorik33/TgBot/internal/models"

type SubscriptionService interface {
	Subscribe(userID int64, token string) error
	GetUserSubscriptions(userID int64) ([]models.Subscription, error)
	Unsubscribe(userID int64, token string) error
	GetAllSubscriptions() ([]models.Subscription, error)
}

type WalletService interface {
	GetWallet(userID int64) ([]models.Portfolio, error)
	AddCryptoToWallet(p *models.Portfolio) error
	DeleteCryptoFromWallet(id int, userID int64) error
}

type CryptoService interface {
	GetCryptoPrice(symbol string) (*models.CryptoInfo, error)
}
