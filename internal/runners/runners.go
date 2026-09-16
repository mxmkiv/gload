package runners

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/mxmkiv/gload/internal/config"
	"github.com/mxmkiv/gload/internal/metrics"
)

type WorkerPool struct {
	config         *config.Config
	HTTPClient     *http.Client
	MetricsChannel chan<- metrics.Metrics
	wg             sync.WaitGroup
}

func NewWorkerPool(cfg *config.Config, HTTPClient *http.Client, metricsChannel chan metrics.Metrics) *WorkerPool {
	return &WorkerPool{
		config:         cfg,
		HTTPClient:     HTTPClient,
		MetricsChannel: metricsChannel,
		wg:             sync.WaitGroup{},
	}
}

func (w *WorkerPool) Start(ctx context.Context) {

	for range w.config.UVs {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			worker(ctx, w)
		}()
	}

	w.wg.Wait()
	close(w.MetricsChannel)

}

func worker(ctx context.Context, w *WorkerPool) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			startTime := time.Now()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.config.Source, nil)
			if err != nil {
				return
			}

			resp, err := w.HTTPClient.Do(req)

			respTime := time.Since(startTime)

			if err != nil {
				w.MetricsChannel <- metrics.Metrics{StatusCode: 0, Latency: respTime, Error: err}
				continue
			}

			_, bodyErr := io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			w.MetricsChannel <- metrics.Metrics{StatusCode: resp.StatusCode, Latency: respTime, Error: bodyErr}
		}
	}
}
