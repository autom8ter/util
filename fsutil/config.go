package fsutil

import (
	"fmt"
	"github.com/autom8ter/util"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string, envPrefix string) {
	viper.SetConfigFile(cfgFile)
	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func SyncEnvConfig() {
	for _, e := range os.Environ() {
		sp := strings.Split(e, "=")
		viper.SetDefault(strings.ToLower(sp[0]), sp[1])
	}
	for k, v := range viper.AllSettings() {
		val, ok := v.(string)
		if ok {
			if err := os.Setenv(k, val); err != nil {
				fmt.Println("failed to bind config to env variable", err.Error())
				os.Exit(1)
			}
		}
	}
}

func RenderFromConfig(s string) string {
	return util.Render(s, viper.AllSettings())
}

func YamlFromConfig() []byte {
	bits, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		fmt.Println("failed to unmarshal current settings to yaml", err.Error())
		os.Exit(1)
	}
	return bits
}

func JsonFromConfig() []byte {
	return util.ToPrettyJson(viper.AllSettings())
}
