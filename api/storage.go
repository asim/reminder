package api

import (
	"os"
	"path/filepath"
)

var ReminderDir = func() string {
	dir := os.ExpandEnv("$HOME/.reminder")
	_ = os.MkdirAll(dir, 0700)
	return dir
}()

func ReminderPath(filename string) string {
	return filepath.Join(ReminderDir, filename)
}
