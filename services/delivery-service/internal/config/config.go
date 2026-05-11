package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort int `mapstructure:"grpc_port"`
	Db       Database
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func bindEnvRecursive(viperInstance *viper.Viper, prefix string, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = strings.ToLower(field.Name[:1]) + field.Name[1:]
		}

		fieldPath := prefix
		if prefix != "" {
			fieldPath = prefix + "." + tag
		} else {
			fieldPath = tag
		}

		if field.Type.Kind() == reflect.Struct {
			if err := bindEnvRecursive(viperInstance, fieldPath, val.Field(i)); err != nil {
				return err
			}
		} else {
			envVarName := strings.ToUpper(strings.ReplaceAll(fieldPath, ".", "_"))
			if err := viperInstance.BindEnv(fieldPath, envVarName); err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetDefault("grpc_port", 8000)

	if err := bindEnvRecursive(v, "", reflect.ValueOf(&Config{}).Elem()); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
