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

var cfg = Config{}

// NewConfig initializes the config
func NewConfig(env string) error {
	v := viper.New()

	v.AutomaticEnv()
	// Replace env key
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	pwd, err := os.Getwd()
	if err != nil {
		return err
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
			return fmt.Errorf("read in config failed, %v", err)
		}
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	if env == "" {
		env = "default"
	}

	slog.Info("Runtime information",
		slog.String("env", env),
		slog.String("goarch", runtime.GOARCH),
		slog.String("goos", runtime.GOOS),
		slog.String("version", runtime.Version()),
	)

	return nil

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
