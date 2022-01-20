package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/config"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer  *http.Server
	httpService *Service
	cfg         config.HTTPConf
	logger      *logrus.Logger
}

func NewServer(logger *logrus.Logger, app app.Application, cfg config.HTTPConf) *Server {
	return &Server{
		logger:      logger,
		cfg:         cfg,
		httpService: NewService(app, logger),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/noop", s.httpService.Noop)

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)),
		Handler: NewLogger(mux),
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
