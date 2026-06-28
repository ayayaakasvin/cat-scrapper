package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	imageprocessing "github.com/ayayaakasvin/cat-scrapper/internal/api-server/libs/image_processing"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/wpn"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func (h *Handlers) SaveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Count *int `json:"count"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil && !errors.Is(err, io.EOF) {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		count := 1
		if req.Count != nil {
			count = *req.Count
		}

		w.WriteHeader(http.StatusAccepted)

		for i := 0; i < count; i++ {
			j := &wpn.Job{
				ID:      ulid.Make().String(),
				Context: context.WithoutCancel(r.Context()),
				Exec: func(ctx context.Context) error {
					img := h.fetchFunc()
					thumbnail, err := imageprocessing.ImageToThumbnail(img.ReaderCloser())
					if err != nil {
						return err
					}

					defer func() { img.Data = nil }()

					uuidString := uuid.NewString()
					fmd := &domain.FileMetaData{
						UUID: uuidString,

						Extension: strings.TrimPrefix(img.ContentType, "image/"),
						MimeType:  img.ContentType,

						Size: int64(len(img.Data)),

						From:      r.RemoteAddr,
						CreatedAt: time.Now(),
					}
					fmd.Filename = uuidString + "." + fmd.Extension
					thumbFilename := uuidString + ".webp"

					// we gotta save file here, but idk should i combine both or try to generic one in order to use saving for further usage
					pathToOriginal, err := h.ifs.SaveImage(img.ReaderCloser(), filepath.Join("original", fmd.Filename))
					if err != nil {
						return err
					}

					pathToThumb, err := h.ifs.SaveImage(io.NopCloser(bytes.NewReader(thumbnail)), filepath.Join("thumb", thumbFilename))
					if err != nil {
						return err
					}

					fmd.Filepath = pathToOriginal
					fmd.Thumbpath = pathToThumb

					if h.fmdr != nil {
						if err := h.fmdr.SaveRecord(context.Background(), fmd.From, fmd); err != nil {
							return err
						}
					}

					return nil
				},
			}

			h.sfn.Submit(j.Exec, r.Context())

		}
	}
}

func (h *Handlers) ServeFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		file, err := h.fmdr.GetByID(r.Context(), id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", file.MimeType)

		http.ServeFile(w, r, file.Filepath)
	}
}

func (h *Handlers) ServeThumbFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		file, err := h.fmdr.GetByID(r.Context(), id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "image/webp")
		w.Header().Set(
			"Cache-Control",
			"public, max-age=31536000, immutable",
		)

		http.ServeFile(w, r, file.Thumbpath)
	}
}
