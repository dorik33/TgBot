package database

import "time"

type Subscription struct {
	ID         int
	UserID     int64
	Token      string
	created_at time.Time
}

// CREATE TABLE IF NOT EXISTS wallet (
//     id SERIAL PRIMARY KEY,
//     user_id BIGINT NOT NULL,
//     token TEXT NOT NULL,
//     Amount FLOAT,
//     Price FLOAT,
//     created_at TIMESTAMP DEFAULT NOW()
// );

type Portfolio struct {
	ID         int
	UserID     int64
	Token      string
	Amount     float64
	Price      float64
	Created_at time.Time
}
