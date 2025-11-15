package handler

import (
	"github.com/Roflan4eg/quiz-api/internal/app/middleware"
)

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required,notblank,min=1,max=5000"`
}

func (r CreateQuestionRequest) Validate() error {
	return middleware.ValidateStruct(r)
}

type CreateAnswerRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Text   string `json:"text" validate:"required,notblank,min=1,max=5000"`
}

func (r CreateAnswerRequest) Validate() error {
	return middleware.ValidateStruct(r)
}
