package logwrapper

import (
	"testing"
)

func TestInvalidArg(t *testing.T) {
	standardLogger := NewLogger()
	defer standardLogger.Sync()
	standardLogger.InvalidArg("nachos")
}
