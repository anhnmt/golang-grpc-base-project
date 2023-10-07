package config

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
