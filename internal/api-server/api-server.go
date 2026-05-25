package apiserver

import (
	"context"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/ayayaakasvin/cat-scrapper/internal/config"
	"github.com/ayayaakasvin/lightmux"
)

type ApiServer struct {
	server *http.Server
	httpcfg *config.Config

	lmux *lightmux.LightMux

	logger *slog.Logger
}

func NewApiServer(
	httpcfg *config.Config,
	logger *slog.Logger,
	) *ApiServer {
	return &ApiServer{

	}
}


func (s *ApiServer) Start(ctx context.Context) error {
	s.setupServer()
	s.setupLightMux()

	go s.printServerStatus(ctx)
	go s.memStatPrint(ctx)
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

	mws := middlewares.NewHTTPMiddlewares(s.logger, *s.corscfg, s.gateawaySecret, s.cache, s.jwtM)
	hndlrs := handlers.NewHTTPHandlers(s.repo, s.cache, s.logger, s.jwtM)

	s.lmux.Use(mws.RecoverMiddleware, mws.GateAwayMiddleware, mws.LoggerMiddleware, mws.CORSMiddleware)

	apiGroup := s.lmux.NewGroup("/api")
	authGroup := apiGroup.ContinueGroup("/auth")

	authGroup.NewRoute("/ping").Handle(http.MethodGet, hndlrs.PingHandler())

	authGroup.NewRoute("/login", mws.RateLimitLoginMiddleware).Handle(http.MethodPost, hndlrs.LogIn())
	authGroup.NewRoute("/register", mws.RateLimitRegisterMiddleware).Handle(http.MethodPost, hndlrs.Register())
	authGroup.NewRoute("/logout", mws.JWTAuthMiddleware).Handle(http.MethodDelete, hndlrs.LogOut())
	authGroup.NewRoute("/refresh").Handle(http.MethodPost, hndlrs.RefreshTheToken())

	authGroup.NewRoute("/public/user").Handle(http.MethodGet, hndlrs.PublicUserInfo())
	authGroup.NewRoute("/me", mws.JWTAuthMiddleware).Handle(http.MethodGet, hndlrs.PrivateUserInfo())

	s.lmux.Mux().HandleFunc("/docs/", httpSwagger.WrapHandler)

	s.logger.Info("LightMux has been set up")
	s.logger.Infof("Available handlers:\n")
	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()
}

func (s *ApiServer) printServerStatus(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Server is alive...")
		case <-ctx.Done():
			return
		}
	}
}

func (s *ApiServer) memStatPrint(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 15)

	select {
	case <-ticker.C:
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		s.logger.Info("Mem stat", "Alloc MiB", m.Alloc/1024/1024)
		time.Sleep(1 * time.Second)
	case <-ctx.Done():
		return
	}
}