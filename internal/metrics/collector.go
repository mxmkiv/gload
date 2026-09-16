package metrics

import (
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

func (c *Collector) Start() {

	for metric := range c.MetricsChannel {
		c.MetricsData = append(c.MetricsData, metric)
	}

}
