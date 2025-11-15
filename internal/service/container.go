package service

import (
	"github.com/Roflan4eg/quiz-api/internal/repository"
)

type Container struct {
	AnswerService   AnswerService
	QuestionService QuestionService
}

func NewContainer(repositoryContainer *repository.Container) *Container {
	qS := NewQuestionService(repositoryContainer.QuestionRepo)
	aS := NewAnswerService(repositoryContainer.AnswerRepo, repositoryContainer.QuestionRepo)
	return &Container{
		AnswerService:   aS,
		QuestionService: qS,
	}
}
