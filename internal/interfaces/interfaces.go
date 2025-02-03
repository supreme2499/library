package interfaces

import (
	"context"
	"library/internal/model"
)

type StorageRepository interface {
	AddNewSong(ctx context.Context, song model.Song) error
	SearchSongsWithFiltering(ctx context.Context, song model.Song, lines, pages int) ([]model.Song, error)
	LyricsWithPagination(ctx context.Context, song, group string, limit, offset int) (verse string, err error)
	EditSongData(ctx context.Context, song model.Song) error
	DeleteSong(ctx context.Context, name, group string) error
}
