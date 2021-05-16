package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	corsConf "github.com/Eldius/cors-interceptor-go/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Setup(cfgFile string) {
	log.Println("Resolving configuration parameters")
	log.Printf("MQTT_BROKER_HOST: %s\n", os.Getenv("MQTT_BROKER_HOST"))
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mqtt-listener-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mqtt-listener-go")
	}

	SetDefaults()
	corsConf.SetDefaults()
	viper.SetEnvPrefix("mqtt")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func SetDefaults() {
	viper.SetDefault("broker.reconnect", true)
}
