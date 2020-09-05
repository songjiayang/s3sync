package config

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/songjiayang/s3sync/syncer"
)

type Config struct {
	Root            string         `json:"root"`
	ScanWorker      int            `json:"scan_worker"`
	DB              string         `json:"db"`
	SyncConfig      *syncer.Config `json:"s3sync"`
	Trim            bool           `json:"trim"`
	AutoContentType bool           `json:"auto_content_type"`
	Interval        int            `json:"interval"` // sync time durting
}

func LoadFile(file string) (cfg *Config, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &cfg)
	return
}

func (c *Config) Duration() time.Duration {
	interval := c.Interval

	if interval == 0 {
		interval = 30
	}

	return time.Duration(interval) * time.Second
}
