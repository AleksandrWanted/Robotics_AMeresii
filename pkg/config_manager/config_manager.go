package config_manager

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type JobConfig struct {
	Id          uint64 `yaml:"id"`
	Name        string `yaml:"name"`
	Schedule    string `yaml:"schedule"`
	Description string `yaml:"description"`
	Method      string `yaml:"method"`
}

type Config struct {
	Jobs []*JobConfig `yaml:"jobs"`
}

type Manager struct {
	cfg *Config
}

func (m *Manager) Config() *Config {
	return m.cfg
}

func NewManager(cfgPath string) *Manager {
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Errorf("can't read file by path: %v, err: %v", cfgPath, err))
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)

	if err != nil {
		panic(fmt.Errorf("can't parse config: %v", err))
	}

	if valid, cause := isValidConfig(cfg); !valid {
		panic(fmt.Errorf("bad config: %v", cause))
	}

	return &Manager{
		cfg: cfg,
	}
}

func CheckCfgAvailability(cfgPath string) (err error) {
	_, err = os.Stat(cfgPath)
	if err != nil {
		return
	}
	return
}

func isValidConfig(cfg *Config) (bool, string) {
	if cfg.Jobs == nil {
		return false, "field jobs can't be omitted"
	}

	if len(cfg.Jobs) == 0 {
		return false, "jobs not found"
	}

	return true, ""
}
