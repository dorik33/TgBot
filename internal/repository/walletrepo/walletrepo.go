package walletrepo

import (
	"database/sql"
	"log"

	"github.com/dorik33/TgBot/internal/models"
	"github.com/dorik33/TgBot/internal/repository"
)

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) repository.WalletRepository {
	return &walletRepository{db: db}
}

func (repo *walletRepository) GetWallet(userID int64) ([]models.Portfolio, error) {
	rows, err := repo.db.Query("SELECT id, user_id, token, amount, price, created_at FROM portfolio WHERE user_id = $1;", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		if err := rows.Scan(&p.ID, &p.UserID, &p.Token, &p.Amount, &p.Price, &p.Created_at); err != nil {
			log.Printf("Error scanning portfolio row: %v", err)
			continue
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}

func (repo *walletRepository) AddCrypto(info *models.Portfolio) error {
	_, err := repo.db.Exec(
		"INSERT INTO portfolio (user_id, token, amount, price) VALUES ($1, $2, $3, $4);",
		info.UserID, info.Token, info.Amount, info.Price,
	)
	return err
}

func (repo *walletRepository) DeleteCrypto(id int, userID int64) error {
	_, err := repo.db.Exec("DELETE FROM portfolio WHERE id = $1 AND user_id = $2;", id, userID)
	return err
}
