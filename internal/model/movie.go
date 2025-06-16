package model

import "time"

type Movie struct {
	ID        int
	GuildID   string
	UserID    string
	Username  string
	Title     string
	Watched   bool
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
