package config

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
