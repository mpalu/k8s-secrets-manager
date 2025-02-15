package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode       string `mapstructure:"mode"`
	KubeConfig string `mapstructure:"kubeconfig"`
	Server     struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
}

func Load() (*Config, error) {
	viper.SetDefault("mode", "cli")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SetConfigFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
}
