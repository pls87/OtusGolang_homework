package grpc

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var requestSeq int64

type ContextKey string

func unaryLoggingInterceptor(l *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		requestID := atomic.AddInt64(&requestSeq, 1)
		remoteAddr := "Unknown"
		if p, e := peer.FromContext(ctx); e {
			remoteAddr = p.Addr.String()
		}
		newCtx := context.WithValue(ctx, ContextKey("request_id"), requestID)

		start := time.Now()
		h, err := handler(newCtx, req)

		l.Infof(`%s [%s] %s request_id: %d request_time: %v error: %v`,
			remoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			info.FullMethod, requestID, time.Since(start), err)

		return h, err
	}
}
