package model

import (
	"time"
)

type Answer struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	QuestionID int       `json:"question_id" validate:"required" gorm:"not null;index"`
	UserID     string    `json:"user_id" validate:"required,uuid" gorm:"type:varchar(36);not null"`
	Text       string    `json:"text" validate:"required,min=1,max=5000" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Answer) TableName() string {
	return "answer"
}
