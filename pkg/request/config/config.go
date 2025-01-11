package config

import (
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/request/shakespeare"
)

type RequestFactoryType string

// Define supported config types here
const (
	ShakespeareRequestFactoryType RequestFactoryType = "Shakespeare"
)

func RequestFactoryFactory(requestType RequestFactoryType) request.RequestFactory {
	switch requestType {
	case ShakespeareRequestFactoryType:
		return shakespeare.ShakespeareRequestFactory
	}
	return nil
}
