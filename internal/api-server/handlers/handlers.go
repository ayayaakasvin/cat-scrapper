// Handlers that serves for main http server, accessed via handlers.Handler struct that contains necessary dependencies
package handlers

import (
	"log/slog"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	saveengine "github.com/ayayaakasvin/cat-scrapper/internal/save-engine"
)

type Handlers struct {
	logger *slog.Logger

	fetchFunc func() *catphotofetch.Image
	saveFunc func(*saveengine.Job, *catphotofetch.Image) (string, error)
}

func NewHTTPHandlers(logger *slog.Logger, fetch func() *catphotofetch.Image, save func(*saveengine.Job, *catphotofetch.Image) (string, error)) *Handlers {
	return &Handlers{
		logger: logger,
		fetchFunc: fetch,
		saveFunc: save,
	}
}