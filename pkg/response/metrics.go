package response

import (
	"time"
)

type Metrics struct {
	Sent    time.Time              `json:"sent"`
	Body    []byte                 `json:"body"`
	Metrics map[string]interface{} `json:"metrics"`
}
