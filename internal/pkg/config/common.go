package config

import (
	"github.com/spf13/viper"
)

func AppName() string {
	return viper.GetString("APP_NAME")
}

func AppPort() int {
	return viper.GetInt("APP_PORT")
}

func PprofEnabled() bool {
	return viper.GetBool("PPROF_ENABLED")
}

func PprofPort() int {
	return viper.GetInt("PPROF_PORT")
}

func LogPayload() bool {
	return viper.GetBool("LOG_PAYLOAD")
}
