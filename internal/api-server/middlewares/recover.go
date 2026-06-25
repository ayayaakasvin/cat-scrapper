package middlewares

import (
	"net/http"
	"runtime/debug"
)

// recover middleware
func (mw *Middlewares) RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				mw.logger.Error("panic recovered", "stack", debug.Stack())

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}
