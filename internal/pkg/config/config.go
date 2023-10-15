package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// New initializes the config
func New(env string) {
	viper.AutomaticEnv()
	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pwd, err := os.Getwd()
	if err != nil {
		log.Panic().Msg(fmt.Sprintf("get current directory failed: %v", err))
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath(pwd)

	envFile := getEnvFile(env)
	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		pe := &fs.PathError{Op: "open", Path: envFile, Err: syscall.ENOENT}
		if ok := errors.As(err, &pe); !ok {
			log.Panic().Msg(fmt.Sprintf("read in config failed: %v", err))
		}
	}

	if env == "" {
		env = "default"
	}

	log.Info().
		Str("env", env).
		Str("app_name", AppName()).
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
