package config

type Config struct {
	// APP
	AppName string `mapstructure:"APP_NAME"`
	AppPort int    `mapstructure:"APP_PORT"`

	// LOG
	LogPayload bool `mapstructure:"LOG_PAYLOAD"`
}

func AppName() string {
	return Default().AppName
}

func AppPort() int {
	return Default().AppPort
}

func LogPayload() bool {
	return Default().LogPayload
}
