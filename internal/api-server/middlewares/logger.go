package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (m *Middlewares) LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		m.logger.Info(requestInfo(r, reqID))
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

func requestInfo(r *http.Request, reqID string) string {
	return fmt.Sprintf(
		"[START]\n\tReqID=%s\n\tMethod=%s\n\tURL=%s\n\tRemoteAddr=%s\n\tUserAgent=%s\n\tHeaders=%v\n",
		reqID,
		r.Method,
		r.URL.String(),
		r.RemoteAddr,
		r.UserAgent(),
		r.Header,
	)
}

// rw implementation for tracking if request was handled successfully, using bool value rw.finished and assigning true if WriteHeader was called.
// Looks like this Req -> Logger -> Handler -> Logger (checks rw.finished value) -> based on it Result shows up.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	finished   bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.finished = true
	rw.ResponseWriter.WriteHeader(code)
}
