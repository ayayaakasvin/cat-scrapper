package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (m *Middlewares) LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		m.logger.Info("[START]",
			slog.String("ReqID", reqID),
			slog.String("Method", r.Method),
			slog.String("URL", r.URL.String()),
			slog.String("RemoteAddr", r.RemoteAddr),
			slog.String("UserAgent", r.UserAgent()),
			slog.Any("Headers", r.Header),
		)
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			result := "failed"
			duration := time.Since(start)

			if rw.finished {
				result = "successfully"
			}

			m.logger.Info("[END]",
				slog.String("ReqID", reqID),
				slog.Int("Status", rw.statusCode),
				slog.Duration("Duration", duration),
				slog.String("Result", result),
			)
		}()

		next.ServeHTTP(rw, r)
	}
}

// rw implementation for tracking if request was handled successfully, using bool value rw.finished and assigning true if WriteHeader was called.
// Looks like this Req -> Logger -> Handler -> Logger (checks rw.finished value) -> based on it Result shows up.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	finished   bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.finished {
		return
	}

	rw.statusCode = code
	rw.finished = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.finished {
		rw.statusCode = http.StatusOK
		rw.finished = true
		rw.ResponseWriter.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}
