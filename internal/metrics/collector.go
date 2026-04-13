package metrics

import (
	"context"
	"fmt"
	"sync"
)

type Collector struct {
	MetricsChannel <-chan Metrics
	Counter        int
	mutex          sync.Mutex
}

func NewCollector(MetricsChannel <-chan Metrics) *Collector {
	return &Collector{
		MetricsChannel: MetricsChannel,
		mutex:          sync.Mutex{},
	}
}

func (c *Collector) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("collector stopped\n")
			return
		case val, ok := <-c.MetricsChannel:
			if !ok {
				return
			}

			c.mutex.Lock()
			c.Counter++
			c.mutex.Unlock()

			fmt.Printf("latency time %v Status code: %v\n", val.Latency, val.StatusCode)
		}
	}

}
