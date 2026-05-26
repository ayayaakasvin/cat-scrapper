package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	saveengine "github.com/ayayaakasvin/cat-scrapper/internal/save-engine"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func (mw *Handlers) SaveHandler() http.HandlerFunc {
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

				j := &saveengine.Job{
					ID:        ulid.Make().String(),
					From:      host,
					ImageUUID: uuid.NewString(),
				}

				img := mw.fetchFunc()

				pathToFile, err := mw.saveFunc(j, img)
				if err != nil {
					mw.logger.Error("Save error", "error", err)
					return
				}

				mw.logger.Info("File saved", "path", pathToFile)
			}()
		}

		wg.Wait()

		w.WriteHeader(http.StatusOK)
	}
}
