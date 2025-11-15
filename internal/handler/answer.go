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

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req CreateAnswerRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = req.Validate()
	if err != nil {
		logger.Error("Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := h.answerService.CreateAnswer(questionID, req.UserID, req.Text)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(answer)
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		middleware.HandleError(w, r, fmt.Errorf("invalid answer ID"))
		return
	}

	answer, err := h.answerService.GetAnswer(id)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 0 {
		middleware.HandleError(w, r, fmt.Errorf("invalid answer ID"))
		return
	}

	err = h.answerService.DeleteAnswer(id)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
