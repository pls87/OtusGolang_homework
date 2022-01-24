package http

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var requestSeq int64

type ContextKey string

type LoggingMiddleware struct {
	handler http.Handler
	logger  *logrus.Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestSeq, 1)
	ctx := context.WithValue(r.Context(), ContextKey("request_id"), requestID)
	wrapped := wrapResponseWriter(w)
	l.handler.ServeHTTP(&wrapped, r.Clone(ctx))
	l.logger.Infof(`%s [%s] %s %s %s %d "%s" request_id: %d request_time: %v`,
		r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method, r.URL, r.Proto, wrapped.Status(), r.UserAgent(), requestID, time.Since(start))
	// log.Printf("%s %s %v request_id: %d", r.Method, r.URL.Path, time.Since(start), requestID)
}

func NewLogger(handlerToWrap http.Handler, logger *logrus.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{handlerToWrap, logger}
}
