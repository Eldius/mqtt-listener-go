package config

import "github.com/spf13/viper"

func GetBrokerHost() string {
	return viper.GetString("broker.host")
}

func GetBrokerPort() int {
	return viper.GetInt("broker.port")
}

func GetBrokerUser() string {
	return viper.GetString("broker.user")
}

func GetBrokerPass() string {
	return viper.GetString("broker.pass")
}

func GetBrokerTopic() string {
	return viper.GetString("broker.topic")
}
