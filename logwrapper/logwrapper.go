package logwrapper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
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
	Stdout *os.File
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*zap.Logger
	*zap.AtomicLevel
}

var logger *StandardLogger
var once *sync.Once = new(sync.Once)

// GetInstance gets a pointer to the shared logger.
func GetInstance(options ...func(*Config)) *StandardLogger {
	once.Do(func() {
		config := Config{
			Env:    PRODUCTION,
			Stdout: os.Stdout,
		}
		config.SetOptions(options...)

		var encoder zapcore.Encoder
		switch config.Env {
		case DEVELOPMENT:
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		default:
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}

		atomicLevel := zap.NewAtomicLevel()
		baseLogger := zap.New(zapcore.NewCore(encoder, zapcore.Lock(config.Stdout), atomicLevel))
		logger = &StandardLogger{baseLogger, &atomicLevel}
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
