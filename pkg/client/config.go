package client

type ClientConfig struct {
	LogLevel    int
	FailOnError bool
}

func NewRecommendedClientConfig() ClientConfig {
	return ClientConfig{
		LogLevel:    1,
		FailOnError: true,
	}
}
