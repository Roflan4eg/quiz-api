package repository

import (
	"errors"
	"fmt"

	"github.com/Roflan4eg/quiz-api/internal/domain"
	"github.com/Roflan4eg/quiz-api/internal/domain/model"
	"gorm.io/gorm"
)

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &answerRepository{db: db}
}

func (r *answerRepository) GetByID(id int) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.First(&answer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}
	return &answer, nil
}

func (r *answerRepository) Create(answer *model.Answer) error {
	err := r.db.Create(answer).Error
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}
	return nil
}

func (r *answerRepository) Delete(id int) error {
	result := r.db.Delete(&model.Answer{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete answer: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrAnswerNotFound
	}
	return nil
}

func (r *answerRepository) GetByQuestionID(questionID int) ([]model.Answer, error) {
	var answers []model.Answer
	err := r.db.Where("question_id = ?", questionID).Find(&answers).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrQuestionNotFound
		}
		return nil, fmt.Errorf("failed to get answers by question id: %w", err)
	}
	return answers, nil
}

func (r *answerRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Answer{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check answer existence: %w", err)
	}
	return count > 0, nil
}
