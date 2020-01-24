package logwrapper

import (
  "go.uber.org/zap"
  "net/http"
	"log"
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

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	baseLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &StandardLogger{baseLogger}
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid argument"}
	invalidArgValueMessage = Event{2, "Invalid value for argument"}
	missingArgMessage      = Event{3, "Missing argument"}
)

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