package application

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"library/internal/http-server/handlers"
	mwLogger "library/internal/http-server/middleware/logger"
)

func SetupRouter(h *handlers.Handler, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/songs", h.AddNewSong)
	router.Get("/songs", h.SearchSongs)
	router.Get("/songs/verse", h.GetVerse)
	router.Put("/songs", h.EditSong)
	router.Delete("/songs", h.DeleteSong)
	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/swagger/doc.json"))) //The url pointing to API definition

	return router
}
