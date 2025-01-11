package openai

import (
	"errors"
)

var (
	ErrInvalidToken = errors.New("token from OpenAI is invalid")
	ErrNoTokens     = errors.New("response from OpenAI has no tokens")
)
