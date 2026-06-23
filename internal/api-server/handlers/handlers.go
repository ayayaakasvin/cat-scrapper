// Handlers that serves for main http server, accessed via handlers.Handler struct that contains necessary dependencies
package handlers

import (
	"log/slog"
	"path/filepath"
	"text/template"

	catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/wpn"
)

type Handlers struct {
	logger *slog.Logger

	fmdr      domain.FileMetaDataRepository
	fetchFunc func() *catphotofetch.Image
	ifs       domain.ImageFileSystem
	sfn       wpn.Sumbitter

	dashboardTmpl *template.Template
}

func NewHTTPHandlers(
	logger *slog.Logger,
	fmdr domain.FileMetaDataRepository,
	fetch func() *catphotofetch.Image,
	ifs domain.ImageFileSystem,
	sfn wpn.Sumbitter,
) *Handlers {
	dashboardHTMLFilePath := filepath.Join(".", "static", "dashboard.html")
	tmpl, err := template.ParseFiles(dashboardHTMLFilePath)
	if err != nil {
		if logger != nil {
			logger.Error("failed to parse dashboard template", "path", dashboardHTMLFilePath, "error", err)
		}
		tmpl = nil
	}

	return &Handlers{
		logger:        logger,
		fmdr:          fmdr,
		fetchFunc:     fetch,
		dashboardTmpl: tmpl,
		ifs:           ifs,
		sfn:           sfn,
	}
}
