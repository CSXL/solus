package config

import "github.com/spf13/viper"

type Config struct {
	*viper.Viper
}

func New() *Config {
	return &Config{Viper: viper.New()}
}

func (c *Config) Read(configName string, configPath string) error {
	c.SetConfigName(configName)
	c.AddConfigPath(configPath)
	c.SetConfigType("yaml")

	err := c.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
