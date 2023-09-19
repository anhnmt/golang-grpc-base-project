package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/viper"
)

// NewConfig initializes the config
func NewConfig(env string) error {
	viper.AutomaticEnv()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")

	envFile := getEnvFile(env)
	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")

	if err := unwrapError(envFile); err != nil {
		return err
	}

	slog.Info("Runtime information",
		slog.String("env", envFile),
		slog.String("goarch", runtime.GOARCH),
		slog.String("goos", runtime.GOOS),
		slog.String("version", runtime.Version()),
	)

	return nil

}

func unwrapError(envFile string) error {
	err := viper.ReadInConfig()

	pe := &fs.PathError{Op: "open", Path: envFile, Err: syscall.ENOENT}
	if ok := errors.As(err, &pe); ok {
		return nil
	}

	return err
}

func getEnvFile(env string) string {
	envFile := fmt.Sprintf(".env.%s", env)
	if checkEnvFileExist(envFile) {
		return envFile
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
