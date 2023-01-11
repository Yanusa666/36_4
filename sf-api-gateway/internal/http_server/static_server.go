package http_server

import (
	"context"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"sf-api-gateway/internal/config"
	"time"
)

type StaticServer struct {
	cfg        *config.Config
	lgr        zerolog.Logger
	httpServer *http.Server
}

func NewStaticServer(cfg *config.Config, lgr zerolog.Logger,
) (*StaticServer, chan error) {
	netListener, err := net.Listen("tcp", cfg.StaticServer.ListenAddress)
	if err != nil {
		lgr.Fatal().Err(err).Msgf("failed to start listener for static server on %s", cfg.StaticServer.ListenAddress)
	}
	lgr.Debug().Msgf("start listener for static server success on %s", cfg.StaticServer.ListenAddress)

	server := &StaticServer{
		cfg: cfg,
		lgr: lgr,
		httpServer: &http.Server{
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       30 * time.Second,
			Handler:           http.FileServer(http.Dir("./webapp")),
		},
	}

	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- server.httpServer.Serve(netListener)
	}()

	return server, listenErrCh
}

func (srv *StaticServer) Shutdown() error {
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
