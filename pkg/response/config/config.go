package config

import (
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/response/empty"
	"github.com/MenD32/Tempest/pkg/response/openai"
	"k8s.io/klog"
)

type ResponseBuilderType string

const (
	OpenAIResponseType ResponseBuilderType = "openai"
	EmptyResponseType  ResponseBuilderType = "empty"
)

func ResponseBuilderFactory(responseType ResponseBuilderType) response.ResponseBuilder {
	klog.Infof("Response Type: %s\n", responseType)
	switch responseType {
	case OpenAIResponseType:
		return openai.OpenAIResponseBuilder
	case EmptyResponseType:
		return empty.EmptyResponseBuilder
	}
	return nil
}
