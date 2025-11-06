package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type Movie struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewMovie(db *sql.DB, logger *slog.Logger) *Movie {
	return &Movie{
		db:     db,
		logger: logger.With(slog.String("component", "repository.movie")),
	}
}

const movieSelectCols = "id, guild_id, user_id, username, title, description, duration, image_url, content_url, created_at, updated_at"

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
		&m.Description,
		&m.ImageURL,
		&m.ContentURL,
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
		r.logger.Debug("failed to query", slog.String("query", query), slog.Any("args", args), slog.Any("err", err))
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		m, err := scanMovie(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movies: %w", err)
		}
		movies = append(movies, *m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
	}

	r.logger.Debug("retrieved movies", slog.Int("count", len(movies)))
	return movies, nil
}

func (r *Movie) AddMovie(movie *model.Movie) (int64, error) {
	query := " INSERT INTO movies (guild_id, user_id, username, title, description, duration, image_url, content_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, movie.GuildID, movie.UserID, movie.Username, movie.Title, movie.Description, movie.ImageURL, movie.ContentURL)
	if err != nil {
		return 0, fmt.Errorf("failed to exec: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last id: %w", err)
	}

	r.logger.Info("added movie", slog.Int64("id", id), slog.String("title", movie.Title), slog.String("username", movie.Username), slog.String("guild_id", movie.GuildID))
	return id, nil
}

func (r *Movie) GetMovieByID(movieID int) (*model.Movie, error) {
	query := fmt.Sprintf("SELECT %s FROM movies WHERE id = ?", movieSelectCols)
	row := r.db.QueryRow(query, movieID)

	movie, err := scanMovie(row)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("movie not found", slog.Int("id", movieID))
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan movie: %w", err)
	}

	r.logger.Debug("retrieved movie", slog.Int64("id", movie.ID), slog.String("title", movie.Title))
	return movie, nil
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

	movie, err := scanMovie(row)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("active movie not found", slog.String("guild_id", guildID))
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan movie: %w", err)
	}

	r.logger.Debug("retrieved active movie", slog.Int64("id", movie.ID), slog.String("title", movie.Title))
	return movie, nil
}

func (r *Movie) UpdateActive(movieID int64, active bool) error {
	query := "UPDATE movies SET active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, active, movieID)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	r.logger.Info("active movie updated", slog.Int64("id", movieID), slog.Bool("active", active))
	return nil
}

func (r *Movie) UpdateWatched(movieID int64, watched bool) error {
	query := "UPDATE movies SET watched = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query, watched, movieID)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	r.logger.Info("watched movie updated", slog.Int64("id", movieID), slog.Bool("active", watched))
	return nil
}

func (r *Movie) DeleteMovie(movieID int64) error {
	query := "DELETE FROM movies WHERE id = ?"
	_, err := r.db.Exec(query, movieID)
	if err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	r.logger.Info("movie deleted", slog.Int64("id", movieID))
	return nil
}
