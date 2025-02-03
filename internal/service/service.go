package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"library/internal/config"
	"library/internal/interfaces"
	"library/internal/lib/logger/sl"
	"library/internal/model"
)

type Service struct {
	log  *slog.Logger
	repo interfaces.StorageRepository
}

func NewService(log *slog.Logger, repo interfaces.StorageRepository) *Service {
	return &Service{log: log, repo: repo}
}

type ResponseApi struct {
	ReleaseDate string `json:"release_date" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Link        string `json:"link" validate:"required"`
}

func (s *Service) AddNewSong(ctx context.Context, song model.Song) error {
	const op = "service.AddNewSong"
	log := s.log.With(slog.String("op", op))
	log.Info("start")

	cfg := config.MustLoad()
	encodedSong := url.QueryEscape(song.Name)
	encodedGroup := url.QueryEscape(song.Group)
	reqURL := fmt.Sprintf("%s/info?group=%s&song=%s", cfg.UrlInfo, encodedGroup, encodedSong)
	log.Info("send api request")
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Debug("fail to send request")
		return err
	}
	log.Info("the response was received successfully")
	defer resp.Body.Close()

	var songDetail ResponseApi
	if err = render.DecodeJSON(resp.Body, &songDetail); err != nil {
		log.Error("failed to decode response body", sl.Err(err))
		return err
	}
	if err = validator.New().Struct(songDetail); err != nil {
		log.Error("invalid response structure", sl.Err(err))
		return err
	}
	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		log.Debug("failed to parse release date", sl.Err(err))
		return err
	}

	song.ReleaseDate = releaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	if err = s.repo.AddNewSong(ctx, song); err != nil {
		log.Error("failed to add new song", sl.Err(err))
		return err
	}
	log.Info("end")
	return nil
}

func (s *Service) SearchSongs(ctx context.Context, song model.Song, year, page string) ([]model.Song, error) {
	const op = "service.SearchSongs"
	log := s.log.With(slog.String("op", op))
	log.Info("start")
	limit := 10

	song.Name = fmt.Sprintf("%%%s%%", song.Name)
	song.Group = fmt.Sprintf("%%%s%%", song.Group)

	offset, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	offset--

	var parsed time.Time
	if year == "" {
		parsed, err = time.Parse("2006-01-02", "2027-01-01")
		if err != nil {
			return nil, err
		}
	} else {
		parsed, err = time.Parse("2006-01-02", year)
		if err != nil {
			return nil, err
		}
	}

	song.ReleaseDate = parsed
	songs, err := s.repo.SearchSongsWithFiltering(ctx, song, limit, offset*10)
	if err != nil {
		log.Error("fail to search songs", sl.Err(err))
		return nil, err
	}
	return songs, nil
}

func (s *Service) GetVerse(ctx context.Context, song, group, lines, page string) (string, error) {
	const op = "service.GetVerse"
	log := s.log.With(slog.String("op", op))
	log.Info("start")
	limit, err := strconv.Atoi(lines)
	if err != nil {
		return "", err
	}
	offset, err := strconv.Atoi(page)
	if err != nil {
		return "", err
	}
	offset--
	verse, err := s.repo.LyricsWithPagination(ctx, song, group, limit, offset)
	if err != nil {
		log.Error("fail to get lyrics", sl.Err(err))
		return "", err
	}
	if verse == "" {
		log.Error("verse == \"\"")
		return "", errors.New("no lyrics")
	}
	log.Info("end")
	return verse, nil
}

func (s *Service) EditSongData(ctx context.Context, song model.Song) error {
	const op = "service.EditSongData"
	log := s.log.With(slog.String("op", op))
	log.Info("start")

	if err := s.repo.EditSongData(ctx, song); err != nil {
		log.Error("failed to edit song", sl.Err(err))
		return err
	}
	log.Info("end")
	return nil
}

func (s *Service) DeleteSong(ctx context.Context, song, group string) error {
	const op = "service.DeleteSong"
	log := s.log.With(slog.String("op", op))

	if err := s.repo.DeleteSong(ctx, song, group); err != nil {
		log.Error("failed to delete song", sl.Err(err))
		return err
	}
	return nil
}
