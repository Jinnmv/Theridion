package ConfigManager

type Config struct {
	Storage StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Type string `json:"type"`
	DSN  string `json:"dsn"`

	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	DBName   string `json:"dbName"`
}
