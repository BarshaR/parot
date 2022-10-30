package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() (config *Config, err error) {
	initDefaults()

	// TODO: init as defaults
	viper.SetConfigType("yaml")
	viper.SetConfigName("parot.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("External configuration file not found")
		} else {
			log.Println("External configuration file found but loaded with errors:", err)
		}
	}

	// TODO: implement config validation to check existence etc

	return &Config{}, err
}

func initDefaults() {
	viper.SetDefault("proxy.port", "8080")
	viper.SetDefault("proxy.record.type", "file")
}

type Config struct {
}

// TODO: privatise these methods and provide a struct to consumers instead
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
