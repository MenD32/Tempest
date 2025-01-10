package client

import "time"

type Metrics struct {
	Sent              time.Time
	Metrics map[string]interface{}
}
