package config

type GithubConfig struct{
	ClientId string `mapstructure:"GITHUB_CLIENTID"`
	ClientSecret string `mapstructure:"GITHUB_SECRET"`
}

func GetGithubConfig(path string) ( *GithubConfig, error ) {
	config := &GithubConfig{}
	if err := LoadConfig(path, config); err != nil {
		return config, err
	}
	return config, nil
}
