package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger        LoggerConf
	StorageSource string `mapstructure:"storage_source"`
	DBStorage     DBStorage
	HTTPServer    HTTPServer
	GRPCServer    GRPCServer
	Rabbit        Rabbit
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type DBStorage struct {
	ConnectionString string `mapstructure:"conn_string"`
}

type HTTPServer struct {
	Host                     string `mapstructure:"host"`
	Port                     string `mapstructure:"port"`
	ReadHeaderTimeoutSeconds int    `mapstructure:"read_header_timeout_sec"`
}

type GRPCServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Rabbit struct {
	ConnectionString string `mapstructure:"conn_string"`
	Exchange         string `mapstructure:"exchange"`
	Queue            string `mapstructure:"queue"`
}

func buildConfig(configFilePath string) (Config, error) {
	var config Config

	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func NewConfig(configFilePath string) (Config, error) {
	return buildConfig(configFilePath)
}
