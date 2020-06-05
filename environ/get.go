package environ

import (
	"os"
	"path/filepath"
)

// String returns the environment variable value as a string
//
// If the environment variable isn't set or empty, it will return defaultValue
func String(name string, defaultValue string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}
	return defaultValue
}

// GetBool returns the environment variable value as a boolean
//
// If the environment variable doesn't exists, is empty or contains an invalid value, false is returned
func GetBool(name string) bool {
	val := String(name, "")
	return val == "1" || val == "true" || val == "yes"
}

// RootPath returns the root path to the executable
func RootPath() string {
	exePath, _ := os.Executable()
	exePath, _ = filepath.Abs(exePath)
	return filepath.Dir(exePath)
}
