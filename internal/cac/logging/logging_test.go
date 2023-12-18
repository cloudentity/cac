package logging_test

import (
	"context"
	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"testing"
)

func TestInitLogging(t *testing.T) {
	tcs := []struct {
		config  logging.Configuration
		handler slog.Handler
		level   slog.Level
	}{
		{
			config: logging.Configuration{
				Level:  "debug",
				Format: "text",
			},
			handler: &slog.TextHandler{},
			level:   slog.LevelDebug,
		},
		{
			config: logging.Configuration{
				Level:  "info",
				Format: "json",
			},
			handler: &slog.JSONHandler{},
			level:   slog.LevelInfo,
		},
	}

	for _, tc := range tcs {
		require.NoError(t, logging.InitLogging(tc.config))
		handler := slog.Default().Handler()
		require.IsType(t, tc.handler, handler)
		require.True(t, handler.Enabled(context.Background(), tc.level))
	}
}
