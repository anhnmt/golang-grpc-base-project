package config

import (
	"github.com/spf13/viper"
)

func DatabaseEnabled() bool {
	return viper.GetBool("DATABASE_ENABLED")
}

func DatabaseDebug() bool {
	return viper.GetBool("DATABASE_DEBUG")
}

func DatabaseMigration() bool {
	return viper.GetBool("DATABASE_MIGRATION")
}

func DatabasePgbouncer() bool {
	return viper.GetBool("DATABASE_PGBOUNCER")
}

func DatabaseMaxOpenConns() int {
	return viper.GetInt("DATABASE_MAX_OPEN_CONNS")
}

func DatabaseMaxIdleConns() int {
	return viper.GetInt("DATABASE_MAX_IDLE_CONNS")
}

func DatabaseMaxLifetime() int {
	return viper.GetInt("DATABASE_MAX_LIFETIME")
}

func DatabaseHost() string {
	return viper.GetString("DATABASE_HOST")
}

func DatabasePort() int {
	return viper.GetInt("DATABASE_PORT")
}

func DatabaseUser() string {
	return viper.GetString("DATABASE_USER")
}

func DatabasePassword() string {
	return viper.GetString("DATABASE_PASSWORD")
}

func DatabaseName() string {
	return viper.GetString("DATABASE_NAME")
}
