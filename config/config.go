package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() (config *Config, err error) {
	initDefaults()

	// TODO: init as defaults
	viper.SetConfigType("yaml")
	viper.SetConfigName("parot.yaml")
	viper.AddConfigPath(".")

	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}

	if viper.IsSet("proxy.port") {
		fmt.Println("proxy port value - ", viper.Get("proxy.port"))
	}

	return &Config{}, err
}

func initDefaults() {
	viper.SetDefault("proxy.port", "8080")
	viper.SetDefault("proxy.record.type", "file")
}

type Config struct {
}

func (config *Config) GetStringValue(key string) string {
	return viper.GetString(key)
}

func (config *Config) GetIntValue(key string) int {
	return viper.GetInt(key)
}

func (config *Config) GetFloatValue(key string) float64 {
	return viper.GetFloat64(key)
}

func (config *Config) GetBoolValue(key string) bool {
	return viper.GetBool(key)
}
