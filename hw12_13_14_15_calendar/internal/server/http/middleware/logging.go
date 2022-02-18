package middleware

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

var requestSeq int64

type ContextKey string

type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

func WrapResponseWriter(w http.ResponseWriter) StatusResponseWriter {
	return StatusResponseWriter{ResponseWriter: w}
}

func (rw *StatusResponseWriter) Status() int {
	return rw.status
}

func (rw *StatusResponseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	if rw.status != 0 {
		return
	}
	rw.status = code
}

func (rw *StatusResponseWriter) Write(bts []byte) (n int, err error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(bts)
}

type LoggingMiddleware struct {
	handler http.Handler
	logger  *logrus.Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestSeq, 1)
	ctx := context.WithValue(r.Context(), ContextKey("request_id"), requestID)
	wrapped := WrapResponseWriter(w)
	l.handler.ServeHTTP(&wrapped, r.Clone(ctx))
	l.logger.Infof(`%s [%s] %s %s %s %d "%s" request_id: %d request_time: %v`,
		r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method, r.URL, r.Proto, wrapped.Status(), r.UserAgent(), requestID, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler, logger *logrus.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{handlerToWrap, logger}
}
