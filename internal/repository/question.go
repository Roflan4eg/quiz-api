package repository

import (
	"errors"
	"fmt"

	"github.com/Roflan4eg/quiz-api/internal/domain"
	"github.com/Roflan4eg/quiz-api/internal/domain/model"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) GetAll() ([]model.Question, error) {
	var res []model.Question
	err := r.db.Find(&res).Error
	if err != nil {
		return []model.Question{}, fmt.Errorf("failed to get questions: %w", err)
	}
	return res, nil
}

func (r *questionRepository) GetByID(id int) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("Answers").First(&question, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}
	return &question, nil
}

func (r *questionRepository) Create(question *model.Question) error {
	err := r.db.Create(question).Error
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}
	return nil
}

func (r *questionRepository) Delete(id int) error {
	result := r.db.Delete(&model.Question{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete question: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrQuestionNotFound
	}
	return nil
}

func (r *questionRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Question{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check answer existence: %w", err)
	}
	return count > 0, nil
}
