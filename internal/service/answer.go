package service

import (
	"fmt"

	"github.com/Roflan4eg/quiz-api/internal/domain"
	"github.com/Roflan4eg/quiz-api/internal/domain/model"
	"github.com/Roflan4eg/quiz-api/internal/repository"
)

type answerService struct {
	answerRepo   repository.AnswerRepository
	questionRepo repository.QuestionRepository
}

func NewAnswerService(answerRepo repository.AnswerRepository, questionRepo repository.QuestionRepository) AnswerService {
	return &answerService{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
	}
}

func (s *answerService) GetAnswer(id int) (*model.Answer, error) {
	return s.answerRepo.GetByID(id)
}

func (s *answerService) CreateAnswer(questionID int, userID, text string) (*model.Answer, error) {
	questionExists, err := s.questionRepo.Exists(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check question existence: %w", err)
	}

	if !questionExists {
		return nil, domain.ErrQuestionNotFound
	}

	answer := &model.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	err = s.answerRepo.Create(answer)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %w", err)
	}

	return answer, nil
}

func (s *answerService) DeleteAnswer(id int) error {
	exists, err := s.answerRepo.Exists(id)
	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrAnswerNotFound
	}

	return s.answerRepo.Delete(id)

}

func (s *answerService) GetAnswersByQuestionID(questionID int) ([]model.Answer, error) {
	questionExists, err := s.questionRepo.Exists(questionID)
	if err != nil {
		return nil, err
	}

	if !questionExists {
		return nil, domain.ErrQuestionNotFound
	}

	return s.answerRepo.GetByQuestionID(questionID)
}
