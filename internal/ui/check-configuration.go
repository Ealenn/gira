package ui

import (
	"github.com/Ealenn/gira/internal/log"

	"github.com/Ealenn/gira/internal/configuration"
)

func CheckConfiguration(logger *log.Logger, configuration *configuration.Configuration, profileName string, profile *configuration.Profile) {
	if profile == nil {
		if profileName == "default" {
			logger.Warn("😕 No Gira configuration found. It looks like Gira has never been set up.")
			logger.Fatal("Please run '%s' to configure %s, or visit the documentation at %s.", "gira config", "Gira", "https://github.com/Ealenn/gira")
		} else {
			logger.Fatal("⚠️  %s\nProfile %s doesn't exist", "Unable to load profile configuration", profileName)
		}
	}

	if !configuration.IsValid(profile) {
		logger.Fatal("⚠️  %s\nPlease run the %s command or change your configuration here %s", "Profile \""+profile.Name+"\" configuration is invalid", "gira config --profile "+profile.Name, configuration.Path)
	}
}
