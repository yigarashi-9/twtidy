package model

import "time"

// ID ...
type ID string

// User ...
type User struct {
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// Tweet ...
type Tweet struct {
	CreatedAt time.Time `json:"created_at"`
	ID        ID        `json:"id"`
	Text      string    `json:"text"`
}
