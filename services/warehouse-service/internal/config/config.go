package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort int `mapstructure:"grpc_port"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetDefault("grpc_port", 8000)
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
