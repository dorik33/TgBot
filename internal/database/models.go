package database

import "time"

type Subscription struct {
	ID         int
	UserID     int64
	Token      string
	created_at time.Time
}
