package config

type Config struct {
	// APP
	AppName string `mapstructure:"APP_NAME"`
	AppPort int    `mapstructure:"APP_PORT"`

	// PPROF
	PprofEnabled bool `mapstructure:"PPROF_ENABLED"`
	PprofPort    int  `mapstructure:"PPROF_PORT"`

	// LOG
	LogPayload bool `mapstructure:"LOG_PAYLOAD"`
}

func AppName() string {
	return Default().AppName
}

func AppPort() int {
	return Default().AppPort
}

func PprofEnabled() bool {
	return Default().PprofEnabled
}

func PprofPort() int {
	return Default().PprofPort
}

func LogPayload() bool {
	return Default().LogPayload
}
