package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ayayaakasvin/wpn"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"io"
	"net/http"
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
				Context: r.Context(),
				Exec: func(ctx context.Context) error {
					img := h.fetchFunc()
					img.UUID = uuid.NewString()
					defer func() { img.Data = nil }()

					pathToFile, err := h.ifs.SaveImage(img)
					if err != nil {
						return err
					}

					if h.fmdr != nil {
						if err := h.fmdr.SaveRecord(r.Context(), r.Host, img, pathToFile); err != nil {
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

		http.ServeFile(w, r, file.Filepath)
	}
}
