package config

import (
	"reflect"
	"strings"

	"github.com/Romasmi/s-shop-microservices/internal/utils/string_utils"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Server   Server
}

type Database struct {
	Host     string `mapstructure:"db_host"`
	Post     string `mapstructure:"db_port"`
	Name     string `mapstructure:"db_name"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
}

type Server struct {
	Port uint
}

func bindEnvRecursive(viperInstance *viper.Viper, prefix string, val reflect.Value) error {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			tag = string_utils.FirstCharToLowerCase(field.Name)
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

func bindAllEnvVars(viperInstance *viper.Viper) error {
	return bindEnvRecursive(viperInstance, "", reflect.ValueOf(&Config{}).Elem())
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v2 := viper.New()
	v2.SetConfigName("override")
	v2.SetConfigType("yaml")
	v2.AddConfigPath(configPath)
	if err := v2.ReadInConfig(); err == nil {
		err := v.MergeConfigMap(v2.AllSettings())
		if err != nil {
			return nil, err
		}
	}

	if err := bindAllEnvVars(v); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
