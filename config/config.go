package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// Port
}

func LoadConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("parot.yaml")
	viper.AddConfigPath(".")
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if viper.IsSet("proxy.port") {
		fmt.Println("proxy port value - ", viper.Get("proxy.port"))
	}
}
