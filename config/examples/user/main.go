package main

import (
	"log"
	"sync"
	"time"

	"github.com/zaynjarvis/fyp/config/sdk"
)

type Config struct {
	sync.RWMutex
	Name    string
	Version int
	Desc    string
}

func (cfg *Config) GetName() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.Name
}

func (cfg *Config) GetVersion() int {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.Version
}

func (cfg *Config) GetDesc() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.Desc
}

func main() {
	var (
		cfg            = Config{}
		name           = "default"
		version uint32 = 1
	)

	go func() {
		if err := sdk.GetConfig(&cfg, name, version, "localhost:3700"); err != nil {
			log.Println(err)
		}
	}()
	var lastCfg Config
	for {
		if cfg.GetName() != lastCfg.Name || cfg.GetVersion() != lastCfg.Version || cfg.GetDesc() != lastCfg.Desc {
			log.Printf("name: %v, version: %v, desc: %v", cfg.Name, cfg.Version, cfg.Desc)
			lastCfg.Name = cfg.GetName()
			lastCfg.Version = cfg.GetVersion()
			lastCfg.Desc = cfg.GetDesc()
		}
		time.Sleep(time.Second)
	}
}
