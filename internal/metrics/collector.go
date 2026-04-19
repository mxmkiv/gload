package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/mxmkiv/gload/internal/config"
)

type Collector struct {
	MetricsChannel <-chan Metrics
	MetricsData    []Metrics
	mutex          sync.Mutex
}

func NewCollector(cfg *config.Config, MetricsChannel <-chan Metrics) *Collector {
	return &Collector{
		MetricsChannel: MetricsChannel,
		mutex:          sync.Mutex{},
		MetricsData:    make([]Metrics, 0, cfg.UVs*int(time.Duration(cfg.Time).Seconds())),
	}
}

func (c *Collector) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			//fmt.Printf("collector stopped\n")
			return
		case val, ok := <-c.MetricsChannel:
			if !ok {
				return
			}

			c.MetricsData = append(c.MetricsData, val)

			//fmt.Printf("latency time %v Status code: %v\n", val.Latency, val.StatusCode)
		}
	}

}
