package service

import (
	"github.com/Roflan4eg/quiz-api/internal/domain/model"
	"github.com/Roflan4eg/quiz-api/internal/repository"
)

type questionService struct {
	questionRepo repository.QuestionRepository
}

func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{questionRepo: questionRepo}
}

func (s *questionService) GetAllQuestions() ([]model.Question, error) {
	return s.questionRepo.GetAll()
}

func (s *questionService) GetQuestion(id int) (*model.Question, error) {
	res, err := s.questionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *questionService) CreateQuestion(text string) (*model.Question, error) {
	question := &model.Question{
		Text: text,
	}

	err := s.questionRepo.Create(question)
	return question, err
}

func (s *questionService) DeleteQuestion(id int) error {
	err := s.questionRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *questionService) QuestionExists(id int) (bool, error) {
	return s.questionRepo.Exists(id)
}
