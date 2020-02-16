package logwrapper

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Environment PRODUCTION or DEVELOPMENT
type Environment int

const (
	PRODUCTION Environment = 0 + iota
	DEVELOPMENT
)

//Config for intializing the logger
type Config struct {
	Env    Environment
	Stderr *os.File
	Stdout *os.File
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*zap.Logger
}

var logger *StandardLogger
var once *sync.Once = new(sync.Once)

// GetInstance gets a pointer to the shared logger.
func GetInstance(options ...func(*Config)) *StandardLogger {
	once.Do(func() {
		config := Config{
			Env:    PRODUCTION,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		config.SetOptions(options...)

		// First, define our level-handling logic.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

		errorLog := zapcore.Lock(config.Stderr)
		debugLog := zapcore.Lock(config.Stdout)

		var encoder zapcore.Encoder
		switch config.Env {
		case DEVELOPMENT:
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		default:
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}

		// Join the outputs, encoders, and level-handling functions into
		// zapcore.Cores, then tee the four cores together.
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, errorLog, highPriority),
			zapcore.NewCore(encoder, debugLog, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		baseLogger := zap.New(core)
		logger = &StandardLogger{baseLogger}
	})
	return logger
}

// SetOptions takes one or more option function and applies them in order to the logger.
func (config *Config) SetOptions(options ...func(*Config)) {
	for _, opt := range options {
		opt(config)
	}
}

//ResetConfig resets the `once` block, so that the config can be reset on the next `GetInstance` call.  Should be only used in testing.
func ResetConfig() {
	once = new(sync.Once)
}
