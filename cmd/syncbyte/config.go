package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var cfgFile string
var conf *Config

type Config struct {
	Server string `json:"server" mapstructure:"server"`
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".syncbyte")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	settings := viper.AllSettings()

	conf = &Config{}
	mapstructure.Decode(settings, conf)
}
