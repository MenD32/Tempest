package response

import "time"

type Metrics struct {
	Sent    time.Time              `json:"sent"`
	Metrics map[string]interface{} `json:"metrics"`
}
