package models

import "time"

type Msg struct {
	ID      int       `db:"id"`
	UserID  int64     `db:"user_id"`
	MsgID   int       `db:"msg_id"`
	Updated time.Time `db:"updated_at"`
	Created time.Time `db:"created_at"`
}

var StatusUserGuest = "guest"

type User struct {
	ID        int64     `db:"id"`
	LastName  string    `db:"last_name"`
	FirstName string    `db:"first_name"`
	Status    string    `db:"status"`
	Created   time.Time `db:"created_at"`
	Updated   time.Time `db:"updated_at"`
}

type UserRequest struct {
	ID        int64
	LastName  string
	FirstName string
	Status    string
}
