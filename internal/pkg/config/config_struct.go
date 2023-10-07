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
