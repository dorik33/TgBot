package database

import "database/sql"

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (repo *SubscriptionRepository) GetSubcriptions(id int64) ([]Subscription, error) {
	rows, err := repo.db.Query("SELECT * FROM subscriptions WHERE user_id = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var sub Subscription
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.Token, &sub.created_at)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (repo *SubscriptionRepository) AddSubscription(id int64, token string) error {
	_, err := repo.db.Exec("INSERT INTO subscriptions (user_id, token) VALUES ($1, $2);", id, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SubscriptionRepository) DeleteSubscription(id int64, token string) error {
	_, err := repo.db.Exec("DELETE FROM subscriptions WHERE user_id = $1 AND token = $2;", id, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SubscriptionRepository) GetAllSubs() []Subscription {
	rows, _ := repo.db.Query("SELECT * FROM subscriptions WHERE user_id = $1;")

	defer rows.Close()

	var subs []Subscription
	for rows.Next() {
		var sub Subscription
		rows.Scan(&sub.ID, &sub.UserID, &sub.Token, &sub.created_at)
		subs = append(subs, sub)
	}
	return subs
}
