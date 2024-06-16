package internal

import (
	"os"
	"path/filepath"
	"runtime"
)

func getConfigFilePath() string {
	if f := os.Getenv("XS_CONFIG"); f != "" {
		if filepath.IsAbs(f) {
			return f
		}
		wd, _ := os.Getwd()
		return filepath.Join(wd, f)
	}

	// default
	home := userHomeDir()
	userDataDir := filepath.Join(home, ".xs")
	return filepath.Join(userDataDir, "config.lua")
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func getDebugFlag() bool {
	v := os.Getenv("XS_DEBUG")
	if v == "1" || v == "true" || v == "TRUE" || v == "True" || v == "yes" || v == "YES" || v == "Yes" || v == "on" || v == "ON" || v == "On" {
		return true
	}
	return false
}

func getNoColorFlag() bool {
	v := os.Getenv("XS_NO_COLOR")
	if v == "1" || v == "true" || v == "TRUE" || v == "True" || v == "yes" || v == "YES" || v == "Yes" || v == "on" || v == "ON" || v == "On" {
		return true
	}
	return false
}
