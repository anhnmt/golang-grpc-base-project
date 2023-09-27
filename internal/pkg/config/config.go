package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var defaultConfig atomic.Value

// Default returns the default Config.
func Default() *Config { return defaultConfig.Load().(*Config) }

// SetDefault makes c the default Config.
func SetDefault(c *Config) {
	defaultConfig.Store(c)
}

// New initializes the config
func New(env string) {
	v := viper.New()

	v.AutomaticEnv()
	// Replace env key
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("get current directory failed: %v", err))
	}

	v.AddConfigPath(".")
	v.AddConfigPath(pwd)

	envFile := getEnvFile(env)
	v.SetConfigFile(envFile)
	v.SetConfigType("env")

	err = v.ReadInConfig()
	if err != nil {
		pe := &fs.PathError{Op: "open", Path: envFile, Err: syscall.ENOENT}
		if ok := errors.As(err, &pe); !ok {
			panic(fmt.Errorf("read in config failed: %v", err))
		}
	}

	c := new(Config)
	err = v.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v", err))
	}
	defaultConfig.Store(c)

	if env == "" {
		env = "default"
	}

	log.Info().
		Str("env", env).
		Str("goarch", runtime.GOARCH).
		Str("goos", runtime.GOOS).
		Str("version", runtime.Version()).
		Msg("runtime information")
}

func getEnvFile(env string) string {
	if env != "" {
		envFile := fmt.Sprintf(".env.%s", env)
		if checkEnvFileExist(envFile) {
			return envFile
		}
	}

	return ".env"
}

func checkEnvFileExist(envFile string) bool {
	// Sử dụng hàm Stat() để kiểm tra tệp tồn tại hay không
	_, err := os.Stat(envFile)
	if err != nil || os.IsNotExist(err) {
		return false
	}

	return true
}
