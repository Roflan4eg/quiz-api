package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Roflan4eg/quiz-api/internal/domain"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)
	})
}

var errorStatusMap = map[error]int{
	domain.ErrQuestionNotFound: http.StatusNotFound,
	domain.ErrAnswerNotFound:   http.StatusNotFound,
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	traceID := GetTraceID(r.Context())
	statusCode := http.StatusInternalServerError

	if code, exists := errorStatusMap[err]; exists {
		statusCode = code
	}

	logger.Error("handler error",
		"trace_id", traceID,
		"error", err.Error(),
		"status_code", statusCode,
		"method", r.Method,
		"path", r.URL.Path,
	)

	userMessage := err.Error()
	if statusCode >= 500 {
		userMessage = "Internal server error"
	}

	sendErrorResponse(w, userMessage, statusCode)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode error response",
			"error", err,
		)
	}
}

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return "unknown"
}
