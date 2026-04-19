package ui

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func Progressbar(ctx context.Context, t time.Duration) {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	start := time.Now()

	for {
		select {
		case <-ticker.C:
			since := time.Since(start)
			percent := float64(since) / float64(t)

			const width = 40
			filled := int(percent * width)
			bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
			fmt.Printf("\r [%s] (%.1f%%)", bar, percent*100)
		case <-ctx.Done():
			return
		}
	}

}
