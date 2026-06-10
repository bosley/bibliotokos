package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	programName  = "bibliotokos"
	themeDark    = "dark"
	themeLight   = "light"
	settingsFile = "settings.json"
)

type settings struct {
	Theme string `json:"theme"`
}

func settingsPath() string {
	return filepath.Join(xdg.DataHome, programName, settingsFile)
}

func defaultSettings() settings {
	return settings{Theme: themeDark}
}

func loadSettings() (settings, error) {
	path := settingsPath()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultSettings(), nil
	}
	if err != nil {
		return defaultSettings(), err
	}
	var s settings
	if err := json.Unmarshal(data, &s); err != nil {
		return defaultSettings(), nil
	}
	if s.Theme != themeLight && s.Theme != themeDark {
		s.Theme = themeDark
	}
	return s, nil
}

func saveSettings(s settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath(), data, 0644)
}

func GetTheme() (string, error) {
	s, err := loadSettings()
	if err != nil {
		return themeDark, err
	}
	return s.Theme, nil
}

func SetDarkTheme() error {
	return setTheme(themeDark)
}

func SetLightTheme() error {
	return setTheme(themeLight)
}

func GetInstallAppName() string {
	return programName
}

func setTheme(theme string) error {
	if theme != themeLight && theme != themeDark {
		return fmt.Errorf("invalid theme: %q", theme)
	}
	s, err := loadSettings()
	if err != nil {
		return err
	}
	s.Theme = theme
	return saveSettings(s)
}
