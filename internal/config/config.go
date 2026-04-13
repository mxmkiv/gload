package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Duration time.Duration

type Config struct {
	Source string   `json:"source"`
	UVs    int      `json:"UVs"`
	Time   Duration `json:"time"`
}

func (d *Duration) UnmarshalJSON(data []byte) error {

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		parse, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("invalid duration string %q: %w", s, err)
		}
		*d = Duration(parse)
		return nil
	}

	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*d = Duration(time.Duration(n) * time.Second)
		return nil
	}

	return fmt.Errorf("invalid duration format: expected string like '10s' or number of second")
}

func NewConfig(fileName string) (*Config, error) {

	// open file
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open config file error: %e", err)
	}
	defer f.Close()

	// parse
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		// advanced error handling
		return nil, fmt.Errorf("parse config file error: %w", err)
	}

	//validate
	if cfg.Source == "" {
		return nil, fmt.Errorf("field source cant be empty")
	}
	if cfg.UVs == 0 {
		return nil, fmt.Errorf("field UVs cant be empty")
	}

	return &cfg, nil
}

func (c *Config) ShowAll() {
	fmt.Printf("source %s\n", c.Source)
	fmt.Printf("uvs %d\n", c.UVs)
	fmt.Printf("time %v\n", c.Time)
}
