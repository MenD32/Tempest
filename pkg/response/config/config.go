package config

import (
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/response/openai"
)

type ResponseBuilderType string

const (
	OpenAIResponseType ResponseBuilderType = "openai"
)

func ResponseBuilderFactory(responseType ResponseBuilderType) response.ResponseBuilder {
	switch responseType {
	case OpenAIResponseType:
		return openai.OpenAIResponseBuilder
	}
	return nil
}
