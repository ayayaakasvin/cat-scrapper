package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ayayaakasvin/wpn"
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

		w.WriteHeader(http.StatusAccepted)

		for _, record := range records {
			rec := record

			j := &wpn.Job{
				ID:      ulid.Make().String(),
				Context: r.Context(),
				Exec: func(ctx context.Context) error {
					err := h.ifs.DeleteImage("original", rec.Filename)
					if err != nil {
						return fmt.Errorf("failed to delete original image: %v", err)
					}

					thumbFileName := rec.UUID + ".webp"
					err = h.ifs.DeleteImage("thumb", thumbFileName)
					if err != nil {
						return fmt.Errorf("failed to delete thumb image: %v", err)
					}

					if h.fmdr != nil {
						if err := h.fmdr.DeleteRecord(rec.UUID); err != nil {
							return fmt.Errorf("file deleted, metadata delete failed: %v", err)
						}
					}

					return nil
				},
			}

			h.sfn.Submit(j.Exec, r.Context())
		}
	}
}
