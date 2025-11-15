package service

import (
	model2 "github.com/Roflan4eg/quiz-api/internal/domain/model"
)

type QuestionService interface {
	GetAllQuestions() ([]model2.Question, error)
	GetQuestion(id int) (*model2.Question, error)
	CreateQuestion(text string) (*model2.Question, error)
	DeleteQuestion(id int) error
	QuestionExists(id int) (bool, error)
}

type AnswerService interface {
	GetAnswer(id int) (*model2.Answer, error)
	CreateAnswer(questionID int, userID, text string) (*model2.Answer, error)
	DeleteAnswer(id int) error
	GetAnswersByQuestionID(questionID int) ([]model2.Answer, error)
}
