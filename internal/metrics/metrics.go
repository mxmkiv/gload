package metrics

import "time"

type Metrics struct {
	//URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}
