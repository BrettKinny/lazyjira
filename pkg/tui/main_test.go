package tui

import (
	"os"
	"testing"
)

// Isolate config dir: selectProject saves last_project to the real auth.json.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "lazyjira-tui-test")
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("LAZYJIRA_CONFIG_DIR", dir)
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
