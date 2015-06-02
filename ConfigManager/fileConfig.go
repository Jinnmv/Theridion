package ConfigManager

import (
	"encoding/json"
	"os"
)

type FileConfig struct {
}

func NewFileCfg(fileName string) (*Config, error) {
	cfg := FileConfig{}

	return cfg.loadConfig(fileName)
}

func (fc *FileConfig) loadConfig(fileName string) (*Config, error) {
	c := Config{}

	if len(fileName) == 0 {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
