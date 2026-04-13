package client

import (
	"net/http"

	"github.com/mxmkiv/gload/internal/config"
)

func NewHTTPClient(cfg *config.Config) *http.Client {
	return &http.Client{}
}
