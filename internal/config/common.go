package config

type Config struct {
	// APP
	AppName string `mapstructure:"APP_NAME"`
	AppPort int    `mapstructure:"APP_PORT"`
}

func AppName() string {
	return cfg.AppName
}

func AppPort() int {
	return cfg.AppPort
}
