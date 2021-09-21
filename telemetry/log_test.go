package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestVoidLogger(t *testing.T) {
	logger := &voidLogger{}

	logger.Level()
	logger.SetLevel("none")
	logger.With("key", "value")
	logger.Debug("debug", "key", "value")
	logger.Debugf("debug %s", "this")
	logger.Info("info", "key", "value")
	logger.Infof("info %s", "this")
	logger.Warn("warn", "key", "value")
	logger.Warnf("warn %s", "this")
	logger.Error("error", "key", "value")
	logger.Errorf("error %s", "this")
	logger.Close()
}

func TestZapLogger(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		logger := &zapLogger{
			config: &zap.Config{
				Level: zap.NewAtomicLevel(),
			},
			logger: zap.NewNop().Sugar(),
		}

		logger.SetLevel("none")
		assert.Equal(t, LevelNone, logger.Level())

		logger.With("key", "value")
		logger.Debug("debug", "key", "value")
		logger.Debugf("debug %s", "this")
		logger.Info("info", "key", "value")
		logger.Infof("info %s", "this")
		logger.Warn("warn", "key", "value")
		logger.Warnf("warn %s", "this")
		logger.Error("error", "key", "value")
		logger.Errorf("error %s", "this")
		logger.Close()
	})

	t.Run("Debug", func(t *testing.T) {
		logger := &zapLogger{
			config: &zap.Config{
				Level: zap.NewAtomicLevel(),
			},
			logger: zap.NewNop().Sugar(),
		}

		logger.SetLevel("debug")
		assert.Equal(t, LevelDebug, logger.Level())

		logger.With("key", "value")
		logger.Debug("debug", "key", "value")
		logger.Debugf("debug %s", "this")
		logger.Info("info", "key", "value")
		logger.Infof("info %s", "this")
		logger.Warn("warn", "key", "value")
		logger.Warnf("warn %s", "this")
		logger.Error("error", "key", "value")
		logger.Errorf("error %s", "this")
		logger.Close()
	})

	t.Run("Info", func(t *testing.T) {
		logger := &zapLogger{
			config: &zap.Config{
				Level: zap.NewAtomicLevel(),
			},
			logger: zap.NewNop().Sugar(),
		}

		logger.SetLevel("info")
		assert.Equal(t, LevelInfo, logger.Level())

		logger.With("key", "value")
		logger.Debug("debug", "key", "value")
		logger.Debugf("debug %s", "this")
		logger.Info("info", "key", "value")
		logger.Infof("info %s", "this")
		logger.Warn("warn", "key", "value")
		logger.Warnf("warn %s", "this")
		logger.Error("error", "key", "value")
		logger.Errorf("error %s", "this")
		logger.Close()
	})

	t.Run("Warn", func(t *testing.T) {
		logger := &zapLogger{
			config: &zap.Config{
				Level: zap.NewAtomicLevel(),
			},
			logger: zap.NewNop().Sugar(),
		}

		logger.SetLevel("warn")
		assert.Equal(t, LevelWarn, logger.Level())

		logger.With("key", "value")
		logger.Debug("debug", "key", "value")
		logger.Debugf("debug %s", "this")
		logger.Info("info", "key", "value")
		logger.Infof("info %s", "this")
		logger.Warn("warn", "key", "value")
		logger.Warnf("warn %s", "this")
		logger.Error("error", "key", "value")
		logger.Errorf("error %s", "this")
		logger.Close()
	})

	t.Run("Error", func(t *testing.T) {
		logger := &zapLogger{
			config: &zap.Config{
				Level: zap.NewAtomicLevel(),
			},
			logger: zap.NewNop().Sugar(),
		}

		logger.SetLevel("error")
		assert.Equal(t, LevelError, logger.Level())

		logger.With("key", "value")
		logger.Debug("debug", "key", "value")
		logger.Debugf("debug %s", "this")
		logger.Info("info", "key", "value")
		logger.Infof("info %s", "this")
		logger.Warn("warn", "key", "value")
		logger.Warnf("warn %s", "this")
		logger.Error("error", "key", "value")
		logger.Errorf("error %s", "this")
		logger.Close()
	})
}
