package config_test

import (
	"reflect"
	"testing"

	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/response/config"
	"github.com/MenD32/Tempest/pkg/response/openai"
)

func TestResponseBuilderFactory(t *testing.T) {
	tests := []struct {
		name         string
		responseType config.ResponseBuilderType
		expected     response.ResponseBuilder
	}{
		{
			name:         "OpenAIResponseType",
			responseType: config.OpenAIResponseType,
			expected:     openai.OpenAIResponseBuilder,
		},
		{
			name:         "NonExistentResponseType",
			responseType: "non-existent",
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.ResponseBuilderFactory(tt.responseType)
			if reflect.ValueOf(actual).Pointer() != reflect.ValueOf(tt.expected).Pointer() {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}

}
