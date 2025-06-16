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
