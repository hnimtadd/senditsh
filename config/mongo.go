package config

import (
	"fmt"
	"log"
	"strings"
)

type MongoConfig struct {
	Source        string `mapstructure:"db_source"`
	Username      string `mapstructure:"db_username"`
	Password      string `mapstructure:"db_password"`
	Database      string `mapstructure:"db_database"`
	AuthSource    string `mapstructure:"db_authsource"`
	AuthMechanism string `mapstructure:"db_authmechanism"`
}

func GetMongoConfig(path string) (*MongoConfig, error) {
	config := &MongoConfig{}
	if err := LoadConfig(path, config); err != nil {
		return nil, err
	}
	log.Printf("mongo config: %v\n", config)
	return config, nil
}

func (c MongoConfig) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("Database: %v\n", c.Database))
	str.WriteString(fmt.Sprintf("UserName: %v\n", c.Username))
	str.WriteString(fmt.Sprintf("Password: %v\n", c.Password))
	str.WriteString(fmt.Sprintf("Source: %v\n", c.Source))
	str.WriteString(fmt.Sprintf("AuthSource: %v\n", c.AuthSource))
	str.WriteString(fmt.Sprintf("AuthMechanism: %v\n", c.AuthMechanism))
	return str.String()
}
