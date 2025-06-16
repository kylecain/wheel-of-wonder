package repository

import (
	"database/sql"
	"log/slog"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) Create(movie string) (int64, error) {
	query := " INSERT INTO movies (name) VALUES (?)"
	result, err := r.db.Exec(query, movie)

	if err != nil {
		slog.Error("failed to insert movie", "error", err, "name", movie)
		return 0, err
	}

	id, _ := result.LastInsertId()
	slog.Info("created movie", "id", id, "name", movie)
	return id, nil
}

func (r *MovieRepository) GetAll() ([]string, error) {
	query := "SELECT name FROM movies"
	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("failed to query movies", "error", err)
		return nil, err
	}
	defer rows.Close()

	var movies []string
	for rows.Next() {
		var movie string
		if err := rows.Scan(&movie); err != nil {
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
