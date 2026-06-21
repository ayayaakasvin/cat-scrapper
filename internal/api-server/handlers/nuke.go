package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/oklog/ulid/v2"
)

func (h *Handlers) NukeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := h.fmdr.GetAllRecords(r.Context())
		if err != nil {
			http.Error(w, "database query error", http.StatusInternalServerError)
			h.logger.Error("failed to fetch records", "err", err)
			return
		}

		if records == nil {
			http.Error(w, "no records found", http.StatusNotFound)
			return
		}

		host := r.Host
		resp := make(map[string]string)
		var mu sync.Mutex
		wg := sync.WaitGroup{}

		for _, record := range records {
			rec := record
			wg.Add(1)

			go func() {
				defer wg.Done()

				j := &domain.Job{
					ID:   ulid.Make().String(),
					From: host,
				}

				err := h.ifs.DeleteImage(rec)
				if err != nil {
					mu.Lock()
					resp[rec.UUID] = err.Error()
					mu.Unlock()
					h.logger.Error("failed to delete image", "job", j.ID, "image", rec.UUID, "from", j.From, "err", err)
					return
				}

				if h.fmdr != nil {
					if err := h.fmdr.DeleteRecord(r.Context(), rec.UUID); err != nil {
						mu.Lock()
						resp[rec.UUID] = fmt.Sprintf("file deleted, metadata delete failed: %v", err)
						mu.Unlock()
						h.logger.Error("failed to delete metadata record", "job", j.ID, "image", rec.UUID, "from", j.From, "err", err)
						return
					}
				}

				mu.Lock()
				resp[rec.UUID] = "successfully deleted"
				mu.Unlock()
			}()
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode dashboard list", http.StatusInternalServerError)
			h.logger.Error("json encode error", "err", err.Error())
			return
		}
	}
}
