package models

import "time"

type Event struct {
	ID          string
	Title       string
	Description string
	Start       string
	End         string
	Created     string
	Updated     string
	Status      string
}

type User struct {
	ID        int64     `db:"id"`
	LastName  string    `db:"last_name"`
	FirstName string    `db:"first_name"`
	Created   time.Time `db:"created_at"`
	Updated   time.Time `db:"updated_at"`
}

type UserRequest struct {
	ID        int64
	LastName  string
	FirstName string
}

type Session struct {
	ID      int       `db:"id"`
	UserID  int64     `db:"user_id"`
	State   string    `db:"state_name"`
	Updated time.Time `db:"updated_at"`
}
