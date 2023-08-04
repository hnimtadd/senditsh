package config

import "github.com/spf13/viper"

type SSHConfig struct {
	Host string `mapstructure:"ssh_host"`
	Port string `mapstructure:"ssh_port"`
}

func LoadConfig(path string, v any) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(v); err != nil {
		return err
	}
	return nil
}

func GetSSHConfig(path string) (*SSHConfig, error) {
	config := &SSHConfig{}
	if err := LoadConfig(path, config); err != nil {
		return config, err
	}
	return config, nil
}
