package model

import "time"

type Movie struct {
	ID        int64
	GuildID   string
	UserID    string
	Username  string
	Title     string
	Watched   bool
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
