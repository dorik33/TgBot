package database

import (
	"database/sql"
	"log"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (repo *WalletRepository) AddCrypto(info *Portfolio) error {
	_, err := repo.db.Exec(
		"INSERT INTO portfolio (user_id, token, amount, price) VALUES ($1, $2, $3, $4);",
		info.UserID, info.Token, info.Amount, info.Price,
	)
	return err
}

func (repo *WalletRepository) GetWallet(userID int64) ([]Portfolio, error) {
	rows, err := repo.db.Query("SELECT id, user_id, token, amount, price, created_at FROM portfolio WHERE user_id = $1;", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []Portfolio
	for rows.Next() {
		var p Portfolio
		if err := rows.Scan(&p.ID, &p.UserID, &p.Token, &p.Amount, &p.Price, &p.Created_at); err != nil {
			log.Printf("Error scanning portfolio row: %v", err)
			continue
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}

func (repo *WalletRepository) DeleteCrypto(id int, userID int64) error {
	_, err := repo.db.Exec("DELETE FROM portfolio WHERE id = $1 AND user_id = $2;", id, userID)
	return err
}



