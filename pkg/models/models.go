package models

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
	ID        int64
	LastName  string
	FirstName string
}
