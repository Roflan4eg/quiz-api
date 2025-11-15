package repository

import (
	model2 "github.com/Roflan4eg/quiz-api/internal/domain/model"
)

type QuestionRepository interface {
	GetAll() ([]model2.Question, error)
	GetByID(id int) (*model2.Question, error)
	Create(question *model2.Question) error
	Delete(id int) error
	Exists(id int) (bool, error)
}

type AnswerRepository interface {
	GetByID(id int) (*model2.Answer, error)
	Create(answer *model2.Answer) error
	Delete(id int) error
	GetByQuestionID(questionID int) ([]model2.Answer, error)
	Exists(id int) (bool, error)
}
