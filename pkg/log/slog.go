package log

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

const logTime = "2006-01-02T15:04:05.000"

var slogLevels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// SetupSlog setting slog default logger
func SetupSlog(conf *viper.Viper) {
	logFile := &lumberjack.Logger{
		Filename:   conf.GetString("log.log_file"),
		Compress:   conf.GetBool("log.compress"),
		MaxAge:     conf.GetInt("log.max_age"),
		MaxSize:    conf.GetInt("log.max_size"),
		MaxBackups: conf.GetInt("log.max_backups"),
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		if a.Value.Kind() == slog.KindTime {
			return slog.String(a.Key, a.Value.Time().Format(logTime))
		}
		return a
	}

	level, ok := slogLevels[strings.ToLower(conf.GetString("log.log_level"))]
	if !ok {
		level = slog.LevelInfo
	}

	// * Text Log Format
	logger := slog.New(slog.NewTextHandler(
		multiWriter,
		&slog.HandlerOptions{
			AddSource:   true,
			ReplaceAttr: replace,
			Level:       level,
		},
	))
	slog.SetDefault(logger)
}
