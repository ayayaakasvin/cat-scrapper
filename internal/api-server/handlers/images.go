package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
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

		host := r.Host

		wg := sync.WaitGroup{}

		for i := 0; i < count; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				j := &domain.Job{
					ID:        ulid.Make().String(),
					From:      host,
					ImageUUID: uuid.NewString(),
				}

				img := h.fetchFunc()

				pathToFile, err := h.saveFunc(j, img)
				if err != nil {
					h.logger.Error("Save error", "error", err)
					return
				}

				if h.fmdr != nil {
					if err := h.fmdr.SaveRecord(j, img, pathToFile); err != nil {
						h.logger.Error("Metadata save error", "path", pathToFile, "error", err)
						return
					}
				}

				h.logger.Info("File saved", "path", pathToFile)
			}()
		}

		wg.Wait()

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) ServeFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		file, err := h.fmdr.GetByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, file.Filepath)
	}
}
