package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mxmkiv/gload/internal/client"
	"github.com/mxmkiv/gload/internal/config"
	"github.com/mxmkiv/gload/internal/metrics"
	"github.com/mxmkiv/gload/internal/runners"
)

func main() {

	cfg, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal("config error: ", err)
	}
	cfg.ShowAll()

	metricsChannel := make(chan metrics.Metrics, 100)
	HTTPClient := client.NewHTTPClient(cfg)
	wp := runners.NewWorkerPool(cfg, HTTPClient, metricsChannel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Time))
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	collector := metrics.NewCollector(metricsChannel)
	go func() {
		defer wg.Done()
		collector.Start(ctx)
	}()

	wp.Start(ctx)
	wg.Wait()

	fmt.Printf("reqs: %v", collector.Counter)
}
