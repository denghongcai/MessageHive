package log

import "testing"

func TestConsole(t *testing.T) {
	NewLogger("console", `{"level": "Info"}`)
	Debug("Debug")
	Info("Info")
}
