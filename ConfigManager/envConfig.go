package ConfigManager

import "os"

type envConfig struct {
}

func NewEnvCfg() *Config {
	cfg := envConfig{}

	return cfg.loadConfig()
}

func (ec envConfig) loadConfig() *Config {
	c := Config{}

	c.Storage.Type = os.Getenv("STORAGE_TYPE")
	c.Storage.DSN = os.Getenv("STORAGE_DSN")

	c.Storage.Hostname = os.Getenv("STORAGE_HOSTNAME")
	c.Storage.Port = os.Getenv("STORAGE_PORT")

	c.Storage.Username = os.Getenv("STORAGE_USERNAME")
	c.Storage.Password = os.Getenv("STORAGE_PASSWORD")
	c.Storage.DBName = os.Getenv("STORAGE_DB_NAME")

	return &c
}
