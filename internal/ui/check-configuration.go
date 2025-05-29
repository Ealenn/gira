package ui

import (
	"github.com/Ealenn/gira/internal/log"

	"github.com/Ealenn/gira/internal/configuration"
)

func CheckConfiguration(logger *log.Logger, profile *configuration.Profile, configuration *configuration.Configuration) {
	if profile == nil {
		logger.Fatal("⚠️  %s\nProfile doesn't exist", "Unable to load profile configuration")
	}

	if !configuration.IsValid(profile) {
		logger.Fatal("⚠️  %s\nPlease run the %s command or change your configuration here %s", "Profile \""+profile.Name+"\" configuration is invalid", "gira config --profile "+profile.Name, configuration.Path)
	}
}
