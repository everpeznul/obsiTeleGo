package config

import (
	"encoding/json"
	"os"
	"sync"
)

var configMu sync.Mutex

type ObsidianConfig struct {
	Exchange string   `json:"exchange"`
	Topics   []string `json:"topics"`
}

func InitConfig(path string) error {
	configMu.Lock()
	defer configMu.Unlock()

	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	cfg := ObsidianConfig{
		Exchange: "tg_router",
		Topics:   []string{},
	}

	return SaveConfigAtomically(path, &cfg)
}

func AddTopic(path string, newTopic string) error {
	configMu.Lock()
	defer configMu.Unlock()

	fileData, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		fileData = []byte(`{"exchange": "tg_router", "topics": []}`)
	}

	var cfg ObsidianConfig
	if err := json.Unmarshal(fileData, &cfg); err != nil {
		return err
	}

	for _, t := range cfg.Topics {
		if t == newTopic {
			return nil
		}
	}

	cfg.Topics = append(cfg.Topics, "tg."+newTopic)

	return SaveConfigAtomically(path, &cfg)
}

func SaveConfigAtomically(path string, cfg *ObsidianConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
