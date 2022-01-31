package grpc

import (
	"context"
	"net"
	"strconv"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/bridge"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	eventSrv   *bridge.EventService
	cfg        configs.APIConf
	logger     *logrus.Logger
}

func New(logger *logrus.Logger, app app.Application, cfg configs.APIConf) basic.Server {
	return &Server{
		eventSrv: bridge.NewService(app.Events(), logger),
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	socket, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)))
	if err != nil {
		s.logger.Fatalf("grpc server - failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer(grpcMiddleware.WithUnaryServerChain(
		unaryLoggingInterceptor(s.logger),
	))

	generated.RegisterCalendarServer(s.grpcServer, s.eventSrv)
	return s.grpcServer.Serve(socket)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}
