package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mxmkiv/gload/internal/client"
	"github.com/mxmkiv/gload/internal/config"
	"github.com/mxmkiv/gload/internal/metrics"
	"github.com/mxmkiv/gload/internal/runners"
	"github.com/mxmkiv/gload/internal/ui"
)

func main() {

	fmt.Println(`
    ________  ___       ________  ________  ________
   |\   ____\|\  \     |\   __  \|\   __  \|\   ___ \
   \ \  \___|\ \  \    \ \  \|\  \ \  \|\  \ \  \_|\ \
    \ \  \  __\ \  \    \ \  \\\  \ \   __  \ \  \ \\ \
     \ \  \|\  \ \  \____\ \  \\\  \ \  \ \  \ \  \_\\ \
      \ \_______\ \_______\ \_______\ \__\ \__\ \_______\
       \|_______|\|_______|\|_______|\|__|\|__|\|_______|
  `)

	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("no --config file specified, using default: config.json")
		*configPath = "config.json"
	}

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal("config error: ", err)
	}
	cfg.PrintConfig()

	metricsChannel := make(chan metrics.Metrics, 100)
	HTTPClient := client.NewHTTPClient(cfg)
	wp := runners.NewWorkerPool(cfg, HTTPClient, metricsChannel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Time))
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	collector := metrics.NewCollector(cfg, metricsChannel)
	go func() {
		defer wg.Done()
		collector.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		ui.Progressbar(ctx, time.Duration(cfg.Time))
	}()

	wp.Start(ctx)
	wg.Wait()

	aggregator := metrics.NewAggregator(collector.MetricsData, cfg)
	aggregator.PrintResult()

	//fmt.Printf("reqs: %v", collector.ReqsCounter)
}
