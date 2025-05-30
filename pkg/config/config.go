package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// NewConfig creates a new viper instance
func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	slog.Info("loading config file from", "path", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		slog.Error("loaded config faile", "path", path, "error", err)
		os.Exit(1)
	}
	return conf
}
