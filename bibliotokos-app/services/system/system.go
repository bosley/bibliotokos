package system

import "bibliotokos/platform"

type SystemService struct{}

func (s *SystemService) GetTheme() (string, error) {
	return platform.GetTheme()
}

func (s *SystemService) SetDarkTheme() error {
	return platform.SetDarkTheme()
}

func (s *SystemService) SetLightTheme() error {
	return platform.SetLightTheme()
}
