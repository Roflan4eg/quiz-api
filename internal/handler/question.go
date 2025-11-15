package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Roflan4eg/quiz-api/internal/app/middleware"
	"github.com/Roflan4eg/quiz-api/pkg/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questionService.GetAllQuestions()
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := req.Validate()
	if err != nil {
		logger.Error("Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	question, err := h.questionService.CreateQuestion(req.Text)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		middleware.HandleError(w, r, fmt.Errorf("invalid answer ID"))
		return
	}

	question, err := h.questionService.GetQuestion(id)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		middleware.HandleError(w, r, fmt.Errorf("invalid answer ID"))
		return
	}

	err = h.questionService.DeleteQuestion(id)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
