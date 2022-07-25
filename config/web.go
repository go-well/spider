package config

//Web 参数
type Web struct {
	Addr     string `yaml:"addr" json:"addr"`
	Debug    bool   `yaml:"debug,omitempty" json:"debug,omitempty"`
	Compress bool   `yaml:"compress" json:"compress,omitempty"`
}

var WebDefault = Web{
	Addr:     ":8080",
	Debug:    false,
	Compress: true,
}
