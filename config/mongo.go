package config

type MongoConfig struct {
	Source        string `mapstructure:"db_source"`
	Username      string `mapstructure:"db_username"`
	Password      string `mapstructure:"db_password"`
	Database      string `mapstructure:"db_database"`
	AuthSource    string `mapstructure:"db_authsource"`
	AuthMechanism string `mapstructure:"db_authmechanism"`
	LogLevel      string `mapstructure:"db_loglevel"`
}

func GetMongoConfig(path string) (*MongoConfig, error) {
	config := &MongoConfig{}
	if err := LoadConfig(path, config); err != nil {
		return nil, err
	}
	return config, nil

}
