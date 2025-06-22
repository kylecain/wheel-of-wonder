package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) AddUser(user *model.User) (int64, error) {
	query := `
		INSERT INTO users (user_id, username, preferred_day_of_week, preferred_time_of_day, preferred_timezone)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			username = excluded.username,
			preferred_day_of_week = excluded.preferred_day_of_week,
			preferred_time_of_day = excluded.preferred_time_of_day,
			preferred_timezone = excluded.preferred_timezone;
	`
	result, err := r.db.Exec(
		query,
		user.UserID,
		user.Username,
		user.PreferredDayOfWeek,
		user.PreferredTimeOfDay,
		user.PreferredTimezone,
	)

	if err != nil {
		return 0, fmt.Errorf("AddUser Error: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddUser Error: %v", err)
	}

	slog.Info("created user", "id", id, "user_id", user.UserID, "username", user.Username)
	return id, nil
}

func (r *User) UserByUserId(userID string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, user_id, username, preferred_day_of_week, preferred_time_of_day, preferred_timezone, created_at, updated_at
		FROM users 
		WHERE user_id = ?
	`

	row := r.db.QueryRow(query, userID)

	if err := row.Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.PreferredDayOfWeek,
		&user.PreferredTimeOfDay,
		&user.PreferredTimezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			slog.Info("no user found", "user_id", userID)
			return nil, nil
		}
		return nil, fmt.Errorf("UserByUserId Error: %v", err)
	}
	slog.Info("retrieved user", "id", user.ID, "user_id", user.UserID)
	return &user, nil
}
