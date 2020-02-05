package logwrapper

import (
	"go.uber.org/zap/zapcore"
  "go.uber.org/zap"
  "net/http"
  "log"
  "sync"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*zap.Logger
}

var logger *StandardLogger
var once sync.Once

// GetInstance gets a pointer to the shared logger.
func GetInstance() *StandardLogger {
  once.Do(func() {
    baseLogger, err := zap.NewProduction()
    if err != nil {
      log.Fatalf("can't initialize zap logger: %v", err)
    }
  
    logger = &StandardLogger{baseLogger}
  })
  return logger
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid argument"}
	invalidArgValueMessage = Event{2, "Invalid value for argument"}
	missingArgMessage      = Event{3, "Missing argument"}
)

// PanicError records the error and then throws a Panic
func (l *StandardLogger) PanicError(msg string, err error) {
  l.Panic(msg, 
    zap.Error(err),
  )
}

// InfoError records the error and doesn't panic
func (l *StandardLogger) InfoError(msg string, err error) {
  l.Error(msg,
    zap.Error(err),
  )
}

// InvalidArg is a standard info message
func (l *StandardLogger) InvalidArg(argumentName string) {
	l.Info(invalidArgMessage.message,
		zap.String("name", argumentName),
	)
}

// InvalidArgValue is a standard info message
func (l *StandardLogger) InvalidArgValue(argumentName string, argumentValue string) {
	l.Info(invalidArgValueMessage.message,
		zap.String("name", argumentName),
		zap.String("value", argumentValue),
	)
}

// MissingArg is a standard info message
func (l *StandardLogger) MissingArg(argumentName string) {
	l.Info(missingArgMessage.message,
		zap.String("name", argumentName),
	)
}

//APIRequest is a standard info message
func (l *StandardLogger) APIRequest(r *http.Request) {
  l.Info("API Request",
  zap.String("method", r.Method),
  zap.String("requestURI", r.RequestURI),
  //zap.Int("response code", r.Response.StatusCode),
)
}

//APIResponse is a standard info message
func (l *StandardLogger) APIResponse(r *http.Request, statusCode int) {
  l.Info("API Request",
  zap.String("method", r.Method),
  zap.String("requestURI", r.RequestURI),
  zap.String("response code", http.StatusText(statusCode)),
)
}

func (l *StandardLogger) String(name string, value string) zapcore.Field {
  return zap.String(name, value)
}