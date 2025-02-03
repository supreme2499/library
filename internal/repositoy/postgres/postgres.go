package repository

import (
	"context"
	"library/internal/lib/logger/sl"
	"library/internal/model"
	"log/slog"

	"library/internal/interfaces"
	"library/internal/storage/postgres"
)

type Repo struct {
	postgres *postgres.Storage
	log      *slog.Logger
}

func NewStorage(storage *postgres.Storage, log *slog.Logger) interfaces.StorageRepository {
	return &Repo{postgres: storage, log: log}
}

func (r *Repo) AddNewSong(ctx context.Context, song model.Song) error {
	const op = "Repository.AddNewSong"
	log := r.log.With(slog.String("op", op))
	log.Info("the beginning of adding a new song")

	q := "INSERT INTO songs (song, group_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.postgres.Pool.Exec(ctx, q, song.Name, song.Group, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		log.Error("error inserting a new song", sl.Err(err))
		return err
	}
	log.Info("successful song insertion")
	return nil
}

func (r *Repo) SearchSongsWithFiltering(ctx context.Context, song model.Song, lines, pages int) ([]model.Song, error) {
	const op = "Repository.SearchSongsWithFiltering"
	log := r.log.With(slog.String("op", op))
	log.Info("the beginning of searching songs")

	q1 := "SELECT * FROM songs WHERE song ILIKE $1 AND group_name ILIKE $2 AND release_date <= $3 ORDER BY release_date DESC LIMIT $4 OFFSET $5"

	rows, err := r.postgres.Pool.Query(ctx, q1, song.Name, song.Group, song.ReleaseDate, lines, pages)
	if err != nil {
		log.Error("error getting songs", sl.Err(err))
		return nil, err
	}
	defer rows.Close()
	var songs []model.Song
	for rows.Next() {
		var song model.Song

		err := rows.Scan(&song.Id, &song.Name, &song.Group, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Error("error getting songs", sl.Err(err))
			return nil, err
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		log.Error("error getting songs", sl.Err(err))
		return nil, err
	}
	log.Info("successful songs search")
	return songs, nil
}

func (r *Repo) LyricsWithPagination(ctx context.Context, song, group string, limit, offset int) (verse string, err error) {
	const op = "Repository.LyricsWithPagination"
	log := r.log.With(slog.String("op", op))
	log.Info("the beginning of searching songs")

	q := "SELECT string_agg(verse, E'\n\n') AS verses FROM (SELECT unnest(string_to_array(text, E'\n\n')) AS verse FROM songs WHERE song ILIKE $1 AND group_name ILIKE $2 LIMIT $3 OFFSET $4) AS verses"

	err = r.postgres.Pool.QueryRow(ctx, q, song, group, limit, offset).Scan(&verse)
	if err != nil {
		log.Error("error getting songs", sl.Err(err))
		return "", err
	}
	log.Info("successful songs search")
	return verse, nil
}

func (r *Repo) EditSongData(ctx context.Context, song model.Song) error {
	log := r.log.With(slog.String("op", "Repository.EditSongData"))
	log.Info("the beginning of editing a song")

	q := "UPDATE songs SET release_date = $1, text = $2, link = $3 WHERE song = $4 AND group_name = $5"

	_, err := r.postgres.Pool.Exec(ctx, q, song.ReleaseDate, song.Text, song.Link, song.Name, song.Group)
	if err != nil {
		log.Error("data update error", sl.Err(err))
		return err
	}
	log.Info("successful update")
	return nil
}

func (r *Repo) DeleteSong(ctx context.Context, name, group string) error {
	const op = "Repository.DeleteSong"
	log := r.log.With(slog.String("op", op))
	log.Info("the beginning of deleting a song")

	q := "DELETE FROM songs WHERE song = $1 AND group_name = $2"

	_, err := r.postgres.Pool.Exec(ctx, q, name, group)
	if err != nil {
		log.Error("error deleting song", sl.Err(err))
		return err
	}
	log.Info("successful song deletion")
	return nil
}
