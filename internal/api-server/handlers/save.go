package handlers

import (
	"net/http"

	saveengine "github.com/ayayaakasvin/cat-scrapper/internal/save-engine"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func (mw *Handlers) SaveOnceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j := &saveengine.Job{
			ID:        ulid.Make().String(),
			From:      r.Host,
			ImageUUID: uuid.NewString(),
		}

		img := mw.fetchFunc()

		pathToFile, err := mw.saveFunc(j, img)
		if err != nil {
			mw.logger.Error("Save error", "error", err)
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}

		mw.logger.Info("File saved", "path", pathToFile)
		w.WriteHeader(http.StatusOK)
	}
}
