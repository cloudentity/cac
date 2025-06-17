package logging

import (
	"context"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

const LevelTrace slog.Level = -8

var DefaultLoggingConfig = func() *Configuration {
	return &Configuration{
		Level: "info",
	}
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

func InitLogging(config *Configuration) (err error) {
	var (
		levelRef = new(slog.LevelVar)
		handler  slog.Handler
		opts     = &slog.HandlerOptions{
			Level: levelRef,
		}
		logger *slog.Logger
	)

	if err = levelRef.UnmarshalText([]byte(config.Level)); err != nil && strings.ToUpper(config.Level) == "TRACE" {
		levelRef.Set(LevelTrace)
	} else if err != nil {
		return err
	}

	switch config.Format {
	case LogFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)

	slog.With("level", levelRef.Level()).Debug("Initiated logging")

	return nil
}

func Trace(msg string, args... any) {
	slog.Log(context.TODO(), LevelTrace, msg, args...)
}
