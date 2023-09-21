package config

type Config struct {
	// APP
	AppName string `mapstructure:"APP_NAME"`
	AppPort int    `mapstructure:"APP_PORT"`
}

func AppName() string {
	return Default().AppName
}

func AppPort() int {
	return Default().AppPort
}
