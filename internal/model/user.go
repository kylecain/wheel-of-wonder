package model

import "time"

type User struct {
	ID                 int64
	UserID             string
	Username           string
	PreferredDayOfWeek string
	PreferredTimeOfDay string
	PreferredTimezone  string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
