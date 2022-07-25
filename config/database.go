package config

//Database 参数
type Database struct {
	Type     string `yaml:"type" json:"type"`
	URL      string `yaml:"url" json:"url"`
	Debug    bool   `yaml:"debug" json:"debug,omitempty"`
	LogLevel int    `json:"log_level" json:"log_level"`
	Sync     bool   `yaml:"sync" json:"sync,omitempty"`
}

var DatabaseDefault = Database{
	Type:     "sqlite",
	URL:      "sqlite3.db",
	Debug:    false,
	LogLevel: 4,
	Sync:     false,
}
