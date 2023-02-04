package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type gql_client struct {
	Endpoint string `mapstructure:"endpoint"`
}

type gql_server struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	GqlClient gql_client `mapstructure:"gql_client"`
	GqlServer gql_server `mapstructure:"gql_server"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Init() {
	// Load the config into viper, otherwise create it
	c.loadConfig(false)

	// Unmarshal the json file's contents into the config struct
	viper.Unmarshal(&c)
}

func (c Config) loadConfig(retry bool) {
	// Set up the config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/gonotes/")
	err := viper.ReadInConfig()

	// If it failed to read in the config, create a new blank one
	if err != nil && retry == false {
		log.Printf("Creating blank config file at $HOME/.config/gonotes\n")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		configPath := "/.config/gonotes/"
		configName := "config.json"

		configData, err := os.ReadFile("example_config.json")
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(homeDir+configPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		dir := homeDir + configPath + configName
		err = os.WriteFile(dir, configData, 0777)
		if err != nil {
			panic(err)
		}

		// Attempt to load the new default config
		c.loadConfig(true)
	}
}
