package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(configPath string) {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.WatchConfig()
}
