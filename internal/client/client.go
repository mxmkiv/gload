package client

import (
	"net/http"
	"time"

	"github.com/mxmkiv/gload/internal/config"
)

func NewHTTPClient(cfg *config.Config) *http.Client {

	transport := http.Transport{
		MaxIdleConnsPerHost:   cfg.UVs,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second, // cfg setting
	}

	return &http.Client{Transport: &transport}
}
