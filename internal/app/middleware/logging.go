package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Roflan4eg/quiz-api/pkg/logger"

	"github.com/google/uuid"
)

type contextKey string

const (
	TraceIDKey contextKey = "trace_id"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var traceID string
		traceIDRaw, err := uuid.NewV7()
		if err != nil {
			traceID = ""
		} else {
			traceID = traceIDRaw.String()
		}
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		logger.Debug("request started",
			"trace_id", traceID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		duration := time.Since(start)
		logger.Info("request completed",
			"trace_id", traceID,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrappedWriter.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
