package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	KubeConfig string       `mapstructure:"kubeconfig"`
	Server     ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	return nil
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func SetConfigFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
}
