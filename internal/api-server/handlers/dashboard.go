package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) DashboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h == nil || h.fmdr == nil {
			http.Error(w, "database repository unavailable", http.StatusServiceUnavailable)
			h.logger.Error("file metadata repo error", "err", "fmdr is nil")
			return
		}

		if h.dashboardTmpl == nil {
			http.Error(w, "dashboard template unavailable", http.StatusInternalServerError)
			h.logger.Error("file metadata repo error", "err", "template is nil")
			return
		}

		files, err := h.fmdr.GetAllRecords()
		if err != nil {
			http.Error(w, "failed to load dashboard data", http.StatusInternalServerError)
			h.logger.Error("file metadata fetch error", "err", err.Error())
			return
		}

		if err := h.dashboardTmpl.Execute(w, files); err != nil {
			http.Error(w, "failed to render dashboard", http.StatusInternalServerError)
			h.logger.Error("template execute error", "err", err.Error())
			return
		}
	}
}

func (h *Handlers) DashboardListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h == nil || h.fmdr == nil {
			http.Error(w, "database repository unavailable", http.StatusServiceUnavailable)
			h.logger.Error("file metadata repo error", "err", "fmdr is nil")
			return
		}

		ids, err := h.fmdr.GetAllIDs()
		if err != nil {
			http.Error(w, "database repository error", http.StatusNotFound)
			h.logger.Error("file metadata repo error", "err", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ids); err != nil {
			http.Error(w, "failed to encode dashboard list", http.StatusInternalServerError)
			if h.logger != nil {
				h.logger.Error("json encode error", "err", err.Error())
			}
		}
	}
}
