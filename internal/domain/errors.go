package domain

import "errors"

var (
	ErrAnswerNotFound   = errors.New("answer not found")
	ErrQuestionNotFound = errors.New("question not found")
)
