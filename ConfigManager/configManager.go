package ConfigManager

import (
	"fmt"
)

type ConfigManager struct {
	Sources map[string]*Config
	Config  Config
}

func New() *ConfigManager {
	cm := ConfigManager{}
	cm.Sources = make(map[string]*Config)

	return &cm
}

func (cm *ConfigManager) LoadSource(name string, c *Config) {
	if _, ok := cm.Sources[name]; ok {
		fmt.Printf("[ConfigManager] WARNING: Config %s is already exist, overwrite", name) // TODO rewrite, use logging
	}
	cm.Sources[name] = c
}

func (cm *ConfigManager) mergeConfigs(c *Config) {
	// TODO : implement
}

func (cm *ConfigManager) LoadAll() {

	cm.LoadSource("DEFAULT", NewDefaultCfg())

	cm.LoadSource("ENV", NewEnvCfg())

	fileCfg, err := NewFileCfg("config.json")
	if err != nil {
		// TODO : Handle exception!
		fmt.Errorf("[CONFIG MANAGER] ERROR: when loading config file: %v", err)
	} else {
		cm.LoadSource("FILE", fileCfg)
	}

	cm.LoadSource("FLAGS", NewFlagsCfg())

}

func (cm ConfigManager) String() string {
	// TODO :implement
	return ""
}
