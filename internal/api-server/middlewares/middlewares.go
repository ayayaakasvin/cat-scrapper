package middlewares

import (
	"log/slog"
	"strings"
	"github.com/ayayaakasvin/cat-scrapper/internal/config"
)

type Middlewares struct {
	logger *slog.Logger

	// CORS config
	allowedOrigins   string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool
}

func NewHTTPMiddlewares(logger *slog.Logger, corsCfg *config.CorsConfig) *Middlewares {
	return &Middlewares{
		logger: logger,

		allowedOrigins:   strings.Join(corsCfg.AllowedOrigins, ","),
		allowedMethods:   strings.Join(corsCfg.AllowedMethods, ","),
		allowedHeaders:   strings.Join(corsCfg.AllowedHeaders, ","),
		allowCredentials: corsCfg.AllowedCredentials,
	}
}
