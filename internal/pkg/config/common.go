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

	// DATABASE
	DatabaseEnabled      bool   `mapstructure:"DATABASE_ENABLED"`
	DatabaseDebug        bool   `mapstructure:"DATABASE_DEBUG"`
	DatabaseMigration    bool   `mapstructure:"DATABASE_MIGRATION"`
	DatabasePgbouncer    bool   `mapstructure:"DATABASE_PGBOUNCER"`
	DatabaseMaxOpenConns int    `mapstructure:"DATABASE_MAX_OPEN_CONNS"`
	DatabaseMaxIdleConns int    `mapstructure:"DATABASE_MAX_IDLE_CONNS"`
	DatabaseMaxLifetime  int    `mapstructure:"DATABASE_MAX_LIFETIME"`
	DatabaseHost         string `mapstructure:"DATABASE_HOST"`
	DatabasePort         int    `mapstructure:"DATABASE_PORT"`
	DatabaseUser         string `mapstructure:"DATABASE_USER"`
	DatabasePassword     string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName         string `mapstructure:"DATABASE_NAME"`
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

func DatabaseEnabled() bool {
	return Default().DatabaseEnabled
}

func DatabaseDebug() bool {
	return Default().DatabaseDebug
}

func DatabaseMigration() bool {
	return Default().DatabaseMigration
}

func DatabasePgbouncer() bool {
	return Default().DatabasePgbouncer
}

func DatabaseMaxOpenConns() int {
	return Default().DatabaseMaxOpenConns
}

func DatabaseMaxIdleConns() int {
	return Default().DatabaseMaxIdleConns
}

func DatabaseMaxLifetime() int {
	return Default().DatabaseMaxLifetime
}

func DatabaseHost() string {
	return Default().DatabaseHost
}

func DatabasePort() int {
	return Default().DatabasePort
}

func DatabaseUser() string {
	return Default().DatabaseUser
}

func DatabasePassword() string {
	return Default().DatabasePassword
}

func DatabaseName() string {
	return Default().DatabaseName
}
