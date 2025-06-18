package subrepo

import (
	"database/sql"

	"github.com/dorik33/TgBot/internal/models"
	"github.com/dorik33/TgBot/internal/repository"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) repository.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (repo *subscriptionRepository) GetSubcriptions(id int64) ([]models.Subscription, error) {
	rows, err := repo.db.Query("SELECT * FROM subscriptions WHERE user_id = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.Token, &sub.Created_at)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (repo *subscriptionRepository) AddSubscription(id int64, token string) error {
	_, err := repo.db.Exec("INSERT INTO subscriptions (user_id, token) VALUES ($1, $2);", id, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *subscriptionRepository) DeleteSubscription(id int64, token string) error {
	_, err := repo.db.Exec("DELETE FROM subscriptions WHERE user_id = $1 AND token = $2;", id, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *subscriptionRepository) GetAllSubs() ([]models.Subscription, error) {
	rows, err := repo.db.Query("SELECT * FROM subscriptions;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err = rows.Scan(&sub.ID, &sub.UserID, &sub.Token, &sub.Created_at)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
