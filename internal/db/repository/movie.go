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

const movieSelectCols = "id, guild_id, user_id, username, title, created_at, updated_at"

type scanner interface {
	Scan(dest ...any) error
}

func scanMovie(s scanner) (*model.Movie, error) {
	var m model.Movie
	if err := s.Scan(
		&m.ID,
		&m.GuildID,
		&m.UserID,
		&m.Username,
		&m.Title,
		&m.CreatedAt,
		&m.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Movie) getMovies(query string, args ...any) ([]model.Movie, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("getMovies Error: %v", err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		m, err := scanMovie(rows)
		if err != nil {
			return nil, fmt.Errorf("getMovies Error: %v", err)
		}
		movies = append(movies, *m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getMovies Error: %v", err)
	}
	slog.Info("retrieved movies", "count", len(movies))
	return movies, nil
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
	query := fmt.Sprintf("SELECT %s FROM movies WHERE id = ?", movieSelectCols)
	row := r.db.QueryRow(query, movieID)

	m, err := scanMovie(row)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("no movie found", "id", movieID)
			return nil, nil
		}
		return nil, fmt.Errorf("GetMovieByID Error: %v", err)
	}
	slog.Info("retrieved movie", "id", m.ID, "title", m.Title)
	return m, nil
}

func (r *Movie) GetAll(guildID string) ([]model.Movie, error) {
	query := fmt.Sprintf("SELECT %s FROM movies WHERE guild_id = ? AND watched = 0 AND active = 0", movieSelectCols)
	return r.getMovies(query, guildID)
}

func (r *Movie) GetAllUnwatched(guildID string) ([]model.Movie, error) {
	query := fmt.Sprintf("SELECT %s FROM movies WHERE guild_id = ? AND watched = 0 AND active = 0 ORDER BY updated_at DESC", movieSelectCols)
	return r.getMovies(query, guildID)
}

func (r *Movie) GetAllWatched(guildID string) ([]model.Movie, error) {
	query := fmt.Sprintf("SELECT %s FROM movies WHERE guild_id = ? AND watched = 1 ORDER BY updated_at DESC", movieSelectCols)
	return r.getMovies(query, guildID)
}

func (r *Movie) GetActive(guildID string) (*model.Movie, error) {
	query := fmt.Sprintf("SELECT %s FROM movies WHERE guild_id = ? AND active = 1 LIMIT 1", movieSelectCols)
	row := r.db.QueryRow(query, guildID)

	m, err := scanMovie(row)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("no active movie found", "guild_id", guildID)
			return nil, nil
		}
		return nil, fmt.Errorf("GetActive Error: %v", err)
	}
	slog.Info("retrieved active movie", "id", m.ID, "title", m.Title)
	return m, nil
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

func (r *Movie) DeleteMovie(movieID int64) error {
	query := "DELETE FROM movies WHERE id = ?"
	_, err := r.db.Exec(query, movieID)
	if err != nil {
		return fmt.Errorf("DeleteMovie Error: %v", err)
	}

	slog.Info("deleted movie", "movie_id", movieID)
	return nil
}
