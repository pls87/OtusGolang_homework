package internalhttp

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestSeq int64

type ContextKey string

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestSeq, 1)
	ctx := context.WithValue(r.Context(), ContextKey("request_id"), requestID)
	reqWithContext := r.Clone(ctx)
	l.handler.ServeHTTP(w, reqWithContext)
	log.Printf("%s %s %v request_id: %d", r.Method, r.URL.Path, time.Since(start), requestID)
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
