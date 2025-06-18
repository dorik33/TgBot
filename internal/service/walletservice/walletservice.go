package walletservice

import (
	"fmt"

	"github.com/dorik33/TgBot/internal/models"
	"github.com/dorik33/TgBot/internal/repository"
	"github.com/dorik33/TgBot/internal/service"
)

type walletService struct {
	repo          repository.WalletRepository
	cryptoService service.CryptoService
}


func NewWalletService(repo repository.WalletRepository, cryptoService service.CryptoService) service.WalletService {
	return &walletService{
		repo:          repo,
		cryptoService: cryptoService,
	}
}

func (s *walletService) GetWallet(userID int64) ([]models.Portfolio, error) {
	return s.repo.GetWallet(userID)
}

func (s *walletService) AddCryptoToWallet(p *models.Portfolio) error {
	if p.Price == 0 {
		crypto, err := s.cryptoService.GetCryptoPrice(p.Token)
		if err != nil {
			return fmt.Errorf("invalid token: %s", p.Token)
		}
		p.Price = crypto.PriceUSD
	}
	return s.repo.AddCrypto(p)
}

func (s *walletService) DeleteCryptoFromWallet(id int, userID int64) error {
	return s.repo.DeleteCrypto(id, userID)
}
