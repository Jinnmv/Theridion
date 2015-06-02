package ConfigManager

import "flag"

type FlagsConfig struct {
}

func NewFlagsCfg() *Config {
	cfg := FlagsConfig{}

	return cfg.loadConfig()
}

func (fc *FlagsConfig) loadConfig() *Config {
	c := Config{}

	flag.StringVar(&c.Storage.Type, "storageType", "", "postgres, sqlite, oracle, mysql")
	flag.StringVar(&c.Storage.DSN, "storageDSN", "", "connection string or filename")

	flag.StringVar(&c.Storage.Hostname, "storageHostname", "", "DB hostname")
	flag.StringVar(&c.Storage.Port, "storagePort", "", "DB port")

	flag.StringVar(&c.Storage.Username, "storageUsername", "", "DB connection username")
	flag.StringVar(&c.Storage.Password, "storageUsername", "", "DB connection password")
	flag.StringVar(&c.Storage.DBName, "storageUsername", "", "DB name")

	flag.Parse()

	return &c
}
