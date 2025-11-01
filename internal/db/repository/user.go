package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type User struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUser(db *sql.DB, logger *slog.Logger) *User {
	return &User{
		db:     db,
		logger: logger.With(slog.String("component", "repository.user")),
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
		return 0, fmt.Errorf("failed to exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last id: %w", err)
	}

	r.logger.Info("created user", slog.Int64("id", id), slog.String("user_id", user.UserID), slog.String("username", user.Username))
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
			slog.Debug("user not found", slog.String("user_id", userID))
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	r.logger.Debug("retrieved user", slog.Int64("id", user.ID), slog.String("user_id", user.UserID))
	return &user, nil
}
