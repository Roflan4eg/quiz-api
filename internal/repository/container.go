package repository

import (
	"github.com/Roflan4eg/quiz-api/internal/storage"
)

type Container struct {
	AnswerRepo   AnswerRepository
	QuestionRepo QuestionRepository
}

func NewContainer(storage *storage.Container) *Container {
	answerRepo := NewAnswerRepository(storage.SQL())
	questionRepo := NewQuestionRepository(storage.SQL())
	return &Container{
		AnswerRepo:   answerRepo,
		QuestionRepo: questionRepo,
	}
}
