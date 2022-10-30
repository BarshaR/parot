package config

import (
	"errors"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	initDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("External configuration file not found, using internal defaults")
		} else {
			return nil, err
		}
	}

	var config Config
	// Unmarshal viper configuration to internal Config model
	if err := unmarshalConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Initlaise defaults to retain config value precedence set by viper
func initDefaults() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("parot.yaml")
	viper.AddConfigPath(".")

	setDefaultsFromConfigOptions()
}

type Config struct {
	ProxyPort            string
	ProxyHostname        string
	ProxyPersist         bool
	ProxyPersistFile     bool
	ProxyPersistDb       bool
	ProxyPersistFilePath string
	ProxyLogPath         string
	ProxyLogLevel        string
}

// Unmarshal to config model
func unmarshalConfig(config *Config) error {
	// TODO: Add more stringent validation and think of a much nicer way to do this

	if port := viper.GetInt("proxy.port"); port > 0 {
		config.ProxyPort = strconv.Itoa(port)
	} else {
		return errors.New("invalid proxy port provided")
	}

	if hostname := viper.GetString("proxy.hostname"); hostname != "" {
		config.ProxyHostname = hostname
	} else {
		return errors.New("invalid hostname provided")
	}

	return nil
}

type ConfigOption struct {
	Key           string
	DefaultValue  string
	Description   string
	AllowedValues []string
	Required      bool
}

func setDefaultsFromConfigOptions() {
	for _, item := range configOptions {
		if item.DefaultValue != "" {
			viper.SetDefault(item.Key, item.DefaultValue)
		}
	}
}

func getConfigOption(key string) *ConfigOption {
	for _, item := range configOptions {
		if item.Key == key {
			return &item
		}
	}
	return nil
}

func getConfigOptionDefault(key string) string {
	if item := getConfigOption(key); item != nil {
		return item.DefaultValue
	}
	return ""
}

var configOptions = []ConfigOption{
	{
		Key:          "proxy.hostname",
		DefaultValue: "localhost",
		Description:  "Hostname of proxy instance",
		Required:     true,
	},
	{
		Key:          "proxy.port",
		DefaultValue: "8080",
		Description:  "Port for proxy instance to listen on",
		Required:     true,
	},
	{
		Key:          "proxy.persist",
		DefaultValue: "false",
		Description:  "Enable persistence of proxy requests",
		Required:     false,
	},
	{
		Key:          "proxy.persist",
		DefaultValue: "./db",
		Description:  "Enable persistence of proxy requests",
		Required:     false,
	},
	{
		Key:          "proxy.persist.file",
		DefaultValue: "./db",
		Description:  "Persist proxy requests to file sytem",
		Required:     false,
	},
	{
		Key:          "proxy.persist.file.path",
		DefaultValue: "./db",
		Description:  "Path of file system persistence directory",
		Required:     false,
	},
	{
		Key:         "proxy.persist.database.name",
		Description: "Relational database flavour",
		Required:    true,
	},
	{
		Key:          "proxy.persist.database.name",
		DefaultValue: "parotdb",
		Description:  "Database name",
		Required:     false,
	},
	{
		Key:          "proxy.persist.database.hostname",
		DefaultValue: "localhost",
		Description:  "Hostname of database server",
		Required:     false,
	},
	{
		Key:          "proxy.persist.database.port",
		DefaultValue: "8",
		Description:  "Database name",
		Required:     true,
	},
}
