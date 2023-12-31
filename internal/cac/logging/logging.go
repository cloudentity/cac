package logging

import (
	"golang.org/x/exp/slog"
	"os"
)

var DefaultLoggingConfig = Configuration{
	Level: "info",
}

type Configuration struct {
	Level  string    `json:"level"`
	Format LogFormat `json:"format"`
}

type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

func InitLogging(config Configuration) (err error) {
	var (
		levelRef = new(slog.LevelVar)
		handler  slog.Handler
		opts     = &slog.HandlerOptions{
			AddSource: true,
			Level:     levelRef,
		}
	)

	if err = levelRef.UnmarshalText([]byte(config.Level)); err != nil {
		return err
	}

	switch config.Format {
	case LogFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.With("logger", logger).Debug("Initiated logging")

	return nil
}
