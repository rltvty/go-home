package logwrapper

import (
	"testing"
)

func TestInvalidArg(t *testing.T) {
	standardLogger := GetInstance()
	defer standardLogger.Sync()
	standardLogger.InvalidArg("nachos")
}
