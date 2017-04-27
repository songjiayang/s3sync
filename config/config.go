package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/file_scan/s3sync"
)

type Config struct {
	Root       string         `json:"root"`
	ScanWorker int            `json:"scan_worker"`
	DB         string         `json:"db"`
	SyncConfig *s3sync.Config `json:"s3sync"`
	Trim       bool           `json:"trim"`
}

func LoadFile(file string) (cfg *Config, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &cfg)
	return
}
