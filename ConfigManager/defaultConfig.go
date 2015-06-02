package ConfigManager

type DefaultConfig struct {
}

func NewDefaultCfg() *Config {
	cfg := DefaultConfig{}

	return cfg.loadConfig()
}

func (dc *DefaultConfig) loadConfig() *Config {
	c := Config{}

	c.Storage.Type = "postgres"
	c.Storage.DSN = "user=theridion dbname=theridion password=theridion host=localhost port=5432"

	c.Storage.Hostname = "localhost"
	c.Storage.Port = "5432"

	c.Storage.Username = "theridion"
	c.Storage.Password = "theridion"
	c.Storage.DBName = "theridion"

	return &c
}
