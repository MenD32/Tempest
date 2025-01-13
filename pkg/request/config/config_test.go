package config_test

import (
	"reflect"
	"testing"

	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/request/config"
	"github.com/MenD32/Tempest/pkg/request/shakespeare"
)

func TestRequestFactoryFactory(t *testing.T) {
	tests := []struct {
		name        string
		requestType config.RequestFactoryType
		expected    request.RequestFactory
	}{
		{
			name:        "Shakespeare",
			requestType: config.ShakespeareRequestFactoryType,
			expected:    shakespeare.ShakespeareRequestFactory,
		},
		{
			name:        "Invalid",
			requestType: "Invalid",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.RequestFactoryFactory(tt.requestType)
			if !(reflect.ValueOf(actual) == reflect.ValueOf(tt.expected)) {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
			}
		})
	}

}
