package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	_ "library/docs"
	resp "library/internal/lib/api/response"
	"library/internal/lib/logger/sl"
	"library/internal/model"
	"library/internal/service"
)

const invalid = "invalid request"

type Handler struct {
	service service.Service
	log     slog.Logger
}

func NewHandler(service service.Service, log slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

type RequestSong struct {
	Name        string    `json:"song" example:"Первый миллион" validate:"required"`
	Group       string    `json:"group" example:"Кишлак" validate:"required"`
	ReleaseDate time.Time `json:"release_date" example:"2024-11-29T10:00:00+03:00" validate:"required"`
	Text        string    `json:"text" example:"текст песни" validate:"required"`
	Link        string    `json:"link" example:"https://youtu.be/SdcNXIPP9UY?si=1QIYeuOSlNsDtZaS" validate:"required"`
}

type Request struct {
	Song  string `json:"song" example:"Первый миллион" validate:"required"`
	Group string `json:"group" example:"Кишлак" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

//	SearchSongs
//	@Summary		Получить песни с фильтрацией по всем полям и пагинацией
//	@Description	Метод возвращает список песен, с фильтрацией и пагинацией(количество песен на странице определено в servise.SearchSongs и является константным).
//	@Description	Номер страницы является параметром запроса. Если не указать страницу, то вернётся первые 10 песен подходящие по фильтру,
//	@Description	если же не указать ни одного фильтра, то метод вернёт первые 10 песен.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string		false	"Название песни для фильтра"
//	@Param			group	query		string		false	"Название группы для фильтра"
//	@Param			year	query		string		false	"Год для фильтра"	Format(yyyy-mm-dd)
//	@Param			page	query		string		false	"Номер страницы"
//	@Success		200		{array}		RequestSong	"В случае успеха сервер вернёт массив JSON"
//	@Failure		400		{object}	Response	"{Status: "ERROR", Error: msg}"
//	@Router			/songs [get]
func (h *Handler) SearchSongs(w http.ResponseWriter, r *http.Request) {
	const op = "handler.SearchSongs"
	log := h.log.With(slog.String("op", op))
	log.Info("start")
	ctx := r.Context()

	name := r.URL.Query().Get("name")
	group := r.URL.Query().Get("group")
	year := r.URL.Query().Get("year")
	page := r.URL.Query().Get("page")

	if page == "" {
		page = "1"
	}

	song := model.Song{
		Name:  name,
		Group: group,
	}

	songs, err := h.service.SearchSongs(ctx, song, year, page)
	if err != nil {
		log.Error("error receiving songs")
		errorHandler(log, invalid, err, w, r)
		return
	}
	log.Info("done")
	log.Info("success")
	render.JSON(w, r, songs)
}

// GetVerse
//	@Summary		Получить текст песни с пагинацией по куплетам
//	@Description	Метод вернёт текс песни c пагинацией по куплетам строкой.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string		true	"Название песни для её поиска"
//	@Param			group	query		string		true	"Название группы для её поиска"
//	@Param			verse	query		string		false	"Количество строк которе необходимо вернуть (по умолчанию 1)"
//	@Param			page	query		string		false	"С какой строки возвращать (по умолчанию 1)"
//	@Success		200		{string}	string		"В случае успеха сервер венёт строку"
//	@Failure		400		{object}	Response	"{Status: "ERROR", Error: msg}"
//	@Router			/songs/verse [get]
func (h *Handler) GetVerse(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetVerse"
	log := h.log.With(slog.String("op", op))
	log.Info("start")
	ctx := r.Context()

	song := r.URL.Query().Get("name")
	group := r.URL.Query().Get("group")
	lines := r.URL.Query().Get("verse")
	page := r.URL.Query().Get("page")

	if song == "" || group == "" {
		log.Error("bad request")
		errorHandler(log, invalid, errors.New("bad request"), w, r)
		return
	}
	if lines == "" {
		lines = "1"
	}
	if page == "" {
		page = "1"
	}

	verses, err := h.service.GetVerse(ctx, song, group, lines, page)
	if err != nil {
		log.Error("error receiving verses")
		errorHandler(log, invalid, err, w, r)
		return
	}
	log.Info("done")
	log.Info("success")
	render.JSON(w, r, verses)
}

//	AddNewSong		inserting a new song
//	@Summary		Добавить песню в библиотеку
//	@Description	После получения запроса сервис делает API запрос, если API даёт ответ с дополнительной информацией о песне,
//	@Description	то мы сохраняем данные в Postgres. Если же API возвращает ошибку, то мы ничего не сохраняем и возвращаем ошибку.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Request		true	"Данные необходимые для добавления."
//	@Success		200		{object}	Response	"Ответ сервера, при успешном сохранении {Status: "OK"}"
//	@Failure		400		{object}	Response	"{Status: "ERROR", Error: msg}"
//	@Router			/songs [post]
func (h *Handler) AddNewSong(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.AddNewSong"
	log := h.log.With(slog.String("op", op))
	log.Info("start")
	ctx := r.Context()

	req, err := decodeAndValidate[Request](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	song := model.Song{
		Name:  req.Song,
		Group: req.Group,
	}

	err = h.service.AddNewSong(ctx, song)
	if err != nil {
		errorHandler(log, "failed to add song", err, w, r)
		return
	}
	log.Info("done")
	log.Info("successful song insertion")
	render.JSON(w, r, Response{Status: "OK"})
}

//	EditSong		Edit song data
//	@Summary		Изменить данные песни
//	@Description	Изменяет данные песни. В качестве индикатора использует название песни и группу
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RequestSong	true	"Данные необходимые для изменения данных"
//	@Success		200		{object}	Response	"Ответ в случае успешного изменения {Status: "OK"}"
//	@Failure		400		{object}	Response	"Ответ в случае возникновения ошибки {Status: "ERROR", Error: msg}"
//	@Router			/songs [put]
func (h *Handler) EditSong(w http.ResponseWriter, r *http.Request) {
	const op = "handler.EditSong"
	log := h.log.With(slog.String("op", op))
	log.Info("start")
	ctx := r.Context()

	req, err := decodeAndValidate[RequestSong](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}

	song := model.Song{
		Name:        req.Name,
		Group:       req.Group,
		ReleaseDate: req.ReleaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	if err = h.service.EditSongData(ctx, song); err != nil {
		errorHandler(log, "failed to edit song", err, w, r)
		return
	}
	log.Info("done")
	log.Info("successful song edition")
	render.JSON(w, r, Response{Status: "OK"})
}

//	DeleteSong		remove a song
//	@Summary		Удалить песню из библиотеки
//	@Description	Удаляет песню из библиотеки по названию и группе
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Request		true	"Данные необходимые для удаления песни"
//	@Success		200		{object}	Response	"Ответ в случае успешного удаления {Status: "OK"}"
//	@Failure		400		{object}	Response	"Ответ в случае возникновения ошибки {Status: "ERROR", Error: msg}"
//	@Router			/songs [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteSong"
	log := h.log.With(slog.String("op", op))
	log.Info("start")
	ctx := r.Context()

	req, err := decodeAndValidate[Request](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}

	err = h.service.DeleteSong(ctx, req.Song, req.Group)
	if err != nil {
		errorHandler(log, "failed to delete song", err, w, r)
		return
	}
	log.Info("done")
	log.Info("successful song delete")
	render.JSON(w, r, Response{Status: "OK"})
}

// Вспомогательные функции
func decodeAndValidate[T any](r *http.Request, log slog.Logger) (*T, error) {
	var req T
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		return nil, fmt.Errorf("failed to decode request")
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		return nil, fmt.Errorf("validation error: %w", err)
	}
	return &req, nil
}

func errorHandler(log *slog.Logger, msg string, err error, w http.ResponseWriter, r *http.Request) {
	log.Error(msg, sl.Err(err))
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, resp.Error(err.Error()))

}
