package apiserver

import (
	"context"
	"log/slog"
	"net/http"

	imagepool "github.com/ayayaakasvin/cat-photo-fetch/image-pool"
	"github.com/ayayaakasvin/cat-scrapper/internal/api-server/handlers"
	"github.com/ayayaakasvin/cat-scrapper/internal/api-server/middlewares"
	"github.com/ayayaakasvin/cat-scrapper/internal/config"
	"github.com/ayayaakasvin/cat-scrapper/internal/domain"
	"github.com/ayayaakasvin/lightmux"
)

type ApiServer struct {
	server  *http.Server
	httpcfg *config.HTTPServerConfig
	corscfg *config.CorsConfig
	lmux    *lightmux.LightMux

	pool *imagepool.CatImagePool
	sg   domain.ImageFileSystem
	fmdr domain.FileMetaDataRepository

	logger *slog.Logger
}

func NewApiServer(
	httpcfg *config.HTTPServerConfig,
	corscfg *config.CorsConfig,
	logger *slog.Logger,
	sg domain.ImageFileSystem,
	fmdr domain.FileMetaDataRepository,
	pool *imagepool.CatImagePool,
) *ApiServer {
	return &ApiServer{
		httpcfg: httpcfg,
		corscfg: corscfg,
		logger:  logger,
		sg:      sg,
		fmdr:    fmdr,
		pool:    pool,
	}
}

func (s *ApiServer) Start(ctx context.Context) error {
	s.setupServer()
	s.setupLightMux()

	return func() error {
		s.logger.Info("Server has been started", "port", s.httpcfg.Address)

		return s.lmux.Run(ctx)
	}()
}

func (s *ApiServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setuping server by pointer, so we dont have to return any value
func (s *ApiServer) setupServer() {
	if s.server == nil {
		s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.httpcfg.Address
	s.server.IdleTimeout = s.httpcfg.IdleTimeout
	s.server.ReadTimeout = s.httpcfg.Timeout
	s.server.WriteTimeout = s.httpcfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ApiServer) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middlewares.NewHTTPMiddlewares(s.logger, s.corscfg)
	hndlrs := handlers.NewHTTPHandlers(s.logger, s.fmdr, s.pool.Get, s.sg)

	s.lmux.Use(mws.RecoverMiddleware, mws.LoggerMiddleware, mws.CORSMiddleware)

	apiGroup := s.lmux.NewGroup("/api")

	apiGroup.NewRoute("/ping").Handle(http.MethodGet, hndlrs.PingHandler())
	apiGroup.NewRoute("/images").Handle(http.MethodPost, hndlrs.SaveHandler())
	apiGroup.NewRoute("/dashboard").Handle(http.MethodGet, hndlrs.DashboardHandler())
	apiGroup.NewRoute("/dashboard/list").Handle(http.MethodGet, hndlrs.DashboardListHandler())
	apiGroup.NewRoute("/dashboard/nuke").Handle(http.MethodDelete, hndlrs.NukeHandler())
	apiGroup.NewRoute("/dashboard/appstat").Handle(http.MethodGet, hndlrs.DashboardAppStatHandler())
	apiGroup.NewRoute("/files").Handle(http.MethodGet, hndlrs.ServeFile())

	s.logger.Info("LightMux has been set up")
	s.logger.Info("Available handlers")
	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()
}
