package util

import (
	"fmt"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string, envPrefix string, safeWrite bool) {
	viper.SetConfigFile(cfgFile)
	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	if safeWrite {
		viper.SafeWriteConfig()
	}
}
