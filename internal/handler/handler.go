package handler

import (
	"net/http"

	"github.com/Roflan4eg/quiz-api/internal/app/middleware"
	"github.com/Roflan4eg/quiz-api/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	questionService service.QuestionService
	answerService   service.AnswerService
}

func NewHandler(container *service.Container) *Handler {
	return &Handler{
		questionService: container.QuestionService,
		answerService:   container.AnswerService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ErrorHandlerMiddleware)

	router.HandleFunc("/questions", h.GetQuestions).Methods("GET")
	router.HandleFunc("/questions", h.CreateQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", h.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", h.DeleteQuestion).Methods("DELETE")

	router.HandleFunc("/questions/{id}/answers", h.CreateAnswer).Methods("POST")
	router.HandleFunc("/answers/{id}", h.GetAnswer).Methods("GET")
	router.HandleFunc("/answers/{id}", h.DeleteAnswer).Methods("DELETE")

	return router
}
