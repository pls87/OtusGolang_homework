package grpc

import (
	"context"
	"net"
	"strconv"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/configs"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/basic"
	"github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar/internal/server/grpc/generated"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	eventSrv   *EventService
	cfg        configs.APIConf
	logger     *logrus.Logger
}

func New(logger *logrus.Logger, app app.Application, cfg configs.APIConf) basic.Server {
	return &Server{
		eventSrv: NewService(app.Events(), logger),
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)))
	if err != nil {
		s.logger.Fatalf("grpc server - failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer(grpcMiddleware.WithUnaryServerChain(
		unaryLoggingInterceptor(s.logger),
	))

	generated.RegisterCalendarServer(s.grpcServer, s.eventSrv)
	s.logger.Info("gRPC server starting...")
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}
