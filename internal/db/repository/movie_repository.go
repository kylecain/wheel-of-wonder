package repository

import (
	"database/sql"
	"log/slog"

	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) Create(movie *model.Movie) (int64, error) {
	query := " INSERT INTO movies (guild_id, user_id, username, title) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, movie.GuildID, movie.UserID, movie.Username, movie.Title)

	if err != nil {
		slog.Error("failed to insert movie", "error", err, "name", movie)
		return 0, err
	}

	id, _ := result.LastInsertId()
	slog.Info("created movie", "id", id, "title", movie.Title, "user", movie.Username)
	return id, nil
}

func (r *MovieRepository) GetAll(guildID string) ([]model.Movie, error) {
	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE guild_id = ?"
	rows, err := r.db.Query(query, guildID)
	if err != nil {
		slog.Error("failed to query movies", "error", err)
		return nil, err
	}
	defer rows.Close()

	var movies []model.Movie
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
			slog.Error("failed to scan movie", "error", err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error iterating over movies", "error", err)
		return nil, err
	}

	slog.Info("retrieved all movies", "count", len(movies))
	return movies, nil
}

func (r *MovieRepository) GetActive(guildID string) (*model.Movie, error) {
	query := "SELECT id, guild_id, user_id, username, title, created_at, updated_at FROM movies WHERE guild_id = ? AND active = 1 LIMIT 1"
	row := r.db.QueryRow(query, guildID)

	var movie model.Movie
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
			return nil, err
		}
		slog.Error("failed to scan active movie", "error", err)
		return nil, err
	}

	slog.Info("retrieved active movie", "id", movie.ID, "title", movie.Title)
	return &movie, nil
}

func (r *MovieRepository) UpdateActive(movieID int, active bool) error {
	query := "UPDATE movies SET active = ? WHERE id = ?"
	_, err := r.db.Exec(query, active, movieID)
	if err != nil {
		slog.Error("failed to update movie active status", "error", err, "movie_id", movieID, "active", active)
		return err
	}

	slog.Info("updated movie active status", "movie_id", movieID, "active", active)
	return nil
}

func (r *MovieRepository) UpdateWatched(movieID int, watched bool) error {
	query := "UPDATE movies SET watched = ? WHERE id = ?"
	_, err := r.db.Exec(query, watched, movieID)
	if err != nil {
		slog.Error("failed to update movie watched status", "error", err, "movie_id", movieID, "watched", watched)
		return err
	}

	slog.Info("updated movie watched status", "movie_id", movieID, "watched", watched)
	return nil
}
