package http_server

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"sf-api-gateway/internal/config"
	"sf-api-gateway/internal/http_server/handlers"
	"time"
)

type Server struct {
	cfg        *config.Config
	lgr        zerolog.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, lgr zerolog.Logger, handler *handlers.Handler,
) (*Server, chan error) {
	netListener, err := net.Listen("tcp", cfg.HttpServer.ListenAddress)
	if err != nil {
		lgr.Fatal().Err(err).Msgf("failed to start listener for http server on %s", cfg.HttpServer.ListenAddress)
	}
	lgr.Debug().Msgf("start listener for http server success on %s", cfg.HttpServer.ListenAddress)

	server := &Server{
		cfg: cfg,
		lgr: lgr,
		httpServer: &http.Server{
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       30 * time.Second,
		},
	}

	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	router.Handler(http.MethodOptions, "/", cors.Default().Handler(router))
	router.POST("/AddNewsComment", handler.Middlware(handler.AddNewsComment))
	router.POST("/GetNewsDetail", handler.Middlware(handler.GetNewsDetail))
	router.POST("/GetLastNews", handler.Middlware(handler.GetLastNews))
	router.GET("/news/:count", handler.Middlware(handler.Cors(handler.GetLastNewsForWebApp)))

	server.httpServer.Handler = router

	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- server.httpServer.Serve(netListener)
	}()

	return server, listenErrCh
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.httpServer.Shutdown(ctx)
	if err != nil {
		srv.lgr.Error().Err(err).Msg("http server grace shutdown finished with error")
		return err
	}

	srv.lgr.Debug().Msg("http server grace shutdown success")
	time.Sleep(10 * time.Millisecond)
	return nil
}
