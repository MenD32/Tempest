package client

type Response interface {
	Metrics() Metrics
}
