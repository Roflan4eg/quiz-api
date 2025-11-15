package model

import (
	"time"
)

type Question struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" validate:"required,min=1,max=5000" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

func (Question) TableName() string {
	return "question"
}
