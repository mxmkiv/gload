package metrics

import (
	"fmt"
	"net/http"
	"slices"
	"sort"
	"time"

	"github.com/mxmkiv/gload/internal/config"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

type Aggregator struct {
	avg time.Duration
	min time.Duration
	max time.Duration
	mid time.Duration

	totalReqs    int
	successReqs  int
	errorReqs    int
	errorPercent float64

	statusCodeStat map[int]int

	p50 time.Duration
	p95 time.Duration
	p99 time.Duration

	RPS float64

	duration time.Duration
}

func NewAggregator(m []Metrics, cfg *config.Config) *Aggregator {
	a := &Aggregator{statusCodeStat: make(map[int]int, len(m))}

	if len(m) == 0 {
		return a
	}

	a.duration = time.Duration(cfg.Time)
	a.totalReqs = len(m)

	latency := make([]time.Duration, 0, len(m))
	var totalLatency time.Duration

	for _, req := range m {
		latency = append(latency, req.Latency)
		totalLatency += req.Latency

		if req.StatusCode != http.StatusOK {
			a.errorReqs++
		} else {
			a.successReqs++
		}

		a.statusCodeStat[req.StatusCode]++
	}

	slices.Sort(latency)
	a.min = latency[0]
	a.max = latency[len(m)-1]
	a.mid = latency[len(m)/2]
	a.avg = totalLatency / time.Duration(len(m))
	a.RPS = float64(a.totalReqs) / a.duration.Seconds()

	a.errorPercent = float64(a.errorReqs) / float64(a.totalReqs) * 100

	a.p50 = percentile(latency, 50)
	a.p95 = percentile(latency, 95)
	a.p99 = percentile(latency, 99)

	return a
}

func (a *Aggregator) PrintResult() {
	line := colorGray + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + colorReset

	fmt.Println()
	fmt.Println(line)
	fmt.Printf("  %s%sGLOAD RESULTS%s\n", colorBold, colorCyan, colorReset)
	fmt.Println(line)

	// requests
	fmt.Printf("\n  %sRequests%s\n", colorBold, colorReset)
	fmt.Printf("  ├─ total    %s%d%s\n", colorYellow, a.totalReqs, colorReset)
	fmt.Printf("  ├─ success  %s%d%s  (%.1f%%)\n", colorGreen, a.successReqs, colorReset, 100-a.errorPercent)
	fmt.Printf("  └─ errors   %s%d%s  (%.1f%%)\n", colorRed, a.errorReqs, colorReset, a.errorPercent)

	// latency
	fmt.Printf("\n  %sLatency%s\n", colorBold, colorReset)
	fmt.Printf("  ├─ avg      %s%v%s\n", colorYellow, a.avg, colorReset)
	fmt.Printf("  ├─ min      %s%v%s\n", colorYellow, a.min, colorReset)
	fmt.Printf("  ├─ max      %s%v%s\n", colorYellow, a.max, colorReset)
	fmt.Printf("  └─ mid      %s%v%s\n", colorYellow, a.mid, colorReset)

	// percentiles
	fmt.Printf("\n  %sPercentiles%s\n", colorBold, colorReset)
	fmt.Printf("  ├─ p50      %s%v%s\n", colorYellow, a.p50, colorReset)
	fmt.Printf("  ├─ p95      %s%v%s\n", colorYellow, a.p95, colorReset)
	fmt.Printf("  └─ p99      %s%v%s\n", colorYellow, a.p99, colorReset)

	// throughput
	fmt.Printf("\n  %sThroughput%s\n", colorBold, colorReset)
	fmt.Printf("  └─ RPS      %s%.1f req/s%s\n", colorYellow, a.RPS, colorReset)

	// status codes
	fmt.Printf("\n  %sStatus Codes%s\n", colorBold, colorReset)
	codes := make([]int, 0, len(a.statusCodeStat))
	for code := range a.statusCodeStat {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for i, code := range codes {
		prefix := "├─"
		if i == len(codes)-1 {
			prefix = "└─"
		}
		color := colorGreen
		if code == 0 || code >= 400 {
			color = colorRed
		}
		fmt.Printf("  %s %s%d%s      %d\n", prefix, color, code, colorReset, a.statusCodeStat[code])
	}

	fmt.Println()
	fmt.Println(line)
	fmt.Println()
}

func percentile(s []time.Duration, p float64) time.Duration {
	idx := int(float64(len(s)) * p / 100)
	if idx >= len(s) {
		idx = len(s) - 1
	}
	return s[idx]
}
