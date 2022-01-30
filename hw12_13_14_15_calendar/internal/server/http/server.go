package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	mux2 "github.com/gorilla/mux"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/http/handler"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/http/middleware"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer  *http.Server
	httpService *handler.Service
	cfg         configs.APIConf
	logger      *logrus.Logger
}

func New(logger *logrus.Logger, app app.Application, cfg configs.APIConf) basic.Server {
	return &Server{
		logger:      logger,
		cfg:         cfg,
		httpService: handler.NewService(app, logger),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := mux2.NewRouter()
	mux.HandleFunc("/noop", s.httpService.Noop)

	mux.HandleFunc("/event", s.httpService.Events().Get).Methods("GET")
	mux.HandleFunc("/event", s.httpService.Events().New).Methods("POST")
	mux.HandleFunc("/event", s.httpService.Events().Update).Methods("PUT")

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)),
		Handler: middleware.NewLogger(mux, s.logger),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
