package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/config"
	"github.com/dorik33/TgBot/internal/repository"
	"github.com/dorik33/TgBot/internal/repository/subrepo"
	"github.com/dorik33/TgBot/internal/repository/walletrepo"
)

type Store struct {
	db                     *sql.DB
	SubscriptionRepository repository.SubscriptionRepository
	WalletRepository       repository.WalletRepository
}

func NewConnection(cfg config.Config) (*Store, error) {
	db, err := sql.Open("postgres", cfg.Database.DBCon)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с базой данных: %w", err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с базой данных: %w", pingErr)
	}

	store := Store{
		db: db,
	}
	store.SubscriptionRepository = subrepo.NewSubscriptionRepository(db)
	store.WalletRepository = walletrepo.NewWalletRepository(db)

	return &store, nil
}
func (s *Store) Close() {
	s.db.Close()
}
