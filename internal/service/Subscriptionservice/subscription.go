package subscriptionservice

import (
	"fmt"

	"github.com/dorik33/TgBot/internal/models"
	"github.com/dorik33/TgBot/internal/repository"
	"github.com/dorik33/TgBot/internal/service"
)

type subscriptionService struct {
	repo          repository.SubscriptionRepository
	cryptoService service.CryptoService
}

func NewSubscriptionService(repo repository.SubscriptionRepository, cryptoService service.CryptoService) service.SubscriptionService {
	return &subscriptionService{
		repo:          repo,
		cryptoService: cryptoService,
	}
}

func (s *subscriptionService) Subscribe(userID int64, token string) error {
	_, err := s.cryptoService.GetCryptoPrice(token)
	if err != nil {
		return fmt.Errorf("invalid token: %s", token)
	}
	return s.repo.AddSubscription(userID, token)
}

func (s *subscriptionService) GetUserSubscriptions(userID int64) ([]models.Subscription, error) {
	return s.repo.GetSubcriptions(userID)
}

func (s *subscriptionService) Unsubscribe(userID int64, token string) error {
	return s.repo.DeleteSubscription(userID, token)
}

func (s *subscriptionService) GetAllSubscriptions() ([]models.Subscription, error) {
	return s.repo.GetAllSubs()
}
