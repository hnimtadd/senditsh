package config

type HTTPConfig struct {
	Host string `mapstructure:"http_host"`
	Port string `mapstructure:"http_port"`
}

func GetHTTPConfig(path string) (*HTTPConfig, error) {
	config := &HTTPConfig{}
	if err := LoadConfig(path, config); err != nil {
		return config, err
	}
	return config, nil
}
