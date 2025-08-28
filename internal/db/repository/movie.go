package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type Movie struct {
	db *sql.DB
}

func NewMovie(db *sql.DB) *Movie {
	return &Movie{
		db: db,
	}
}

func (r *Movie) AddMovie(movie *model.Movie) (int64, error) {
	query := " INSERT INTO movies (guild_id, user_id, username, title) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, movie.GuildID, movie.UserID, movie.Username, movie.Title)

	if err != nil {
		return 0, fmt.Errorf("AddMovie Error: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddMovie Error: %v", err)
	}

	slog.Info("created movie", "id", id, "title", movie.Title, "user", movie.Username)
	return id, nil
}

func (r *Movie) GetMovieByID(movieID int) (*model.Movie, error) {
	var movie model.Movie

	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE id = ?"
	row := r.db.QueryRow(query, movieID)

	if err := row.Scan(
		&movie.ID,
		&movie.GuildID,
		&movie.UserID,
		&movie.Username,
		&movie.Title,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			slog.Info("no movie found", "id", movieID)
			return nil, nil
		}
		return nil, fmt.Errorf("GetMovieByID Error: %v", err)
	}
	slog.Info("retrieved movie", "id", movie.ID, "title", movie.Title)
	return &movie, nil
}

func (r *Movie) GetAll(guildID string) ([]model.Movie, error) {
	var movies []model.Movie

	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE guild_id = ? AND watched = 0 AND active = 0"
	rows, err := r.db.Query(query, guildID)
	if err != nil {
		return nil, fmt.Errorf("GetAll Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie model.Movie

		if err := rows.Scan(
			&movie.ID,
			&movie.GuildID,
			&movie.UserID,
			&movie.Username,
			&movie.Title,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetAll Error: %v", err)
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAll Error: %v", err)
	}

	slog.Info("retrieved all movies", "count", len(movies))
	return movies, nil
}

func (r *Movie) GetAllWatched(guildID string) ([]model.Movie, error) {
	var movies []model.Movie

	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE guild_id = ? AND watched = 1 ORDER BY updated_at DESC"
	rows, err := r.db.Query(query, guildID)
	if err != nil {
		return nil, fmt.Errorf("GetAllWatched Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie model.Movie

		if err := rows.Scan(
			&movie.ID,
			&movie.GuildID,
			&movie.UserID,
			&movie.Username,
			&movie.Title,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetAllWatched Error: %v", err)
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllWatched Error: %v", err)
	}

	slog.Info("retrieved watched movies", "count", len(movies))
	return movies, nil
}

func (r *Movie) GetActive(guildID string) (*model.Movie, error) {
	var movie model.Movie

	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE guild_id = ? AND active = 1 LIMIT 1"
	row := r.db.QueryRow(query, guildID)

	if err := row.Scan(
		&movie.ID,
		&movie.GuildID,
		&movie.UserID,
		&movie.Username,
		&movie.Title,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			slog.Info("no active movie found", "guild_id", guildID)
			return nil, nil
		}
		return nil, fmt.Errorf("GetActive Error: %v", err)
	}

	slog.Info("retrieved active movie", "id", movie.ID, "title", movie.Title)
	return &movie, nil
}

func (r *Movie) UpdateActive(movieID int64, active bool) error {
	query := "UPDATE movies SET active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, active, movieID)
	if err != nil {
		return fmt.Errorf("UpdateActive Error: %v", err)
	}

	slog.Info("updated movie active status", "movie_id", movieID, "active", active)
	return nil
}

func (r *Movie) UpdateWatched(movieID int64, watched bool) error {
	query := "UPDATE movies SET watched = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, watched, movieID)
	if err != nil {
		return fmt.Errorf("UpdateWatched Error: %v", err)
	}

	slog.Info("updated movie watched status", "movie_id", movieID, "watched", watched)
	return nil
}
