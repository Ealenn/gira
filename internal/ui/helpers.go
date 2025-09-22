package ui

import (
	"time"

	"github.com/Ealenn/gira/internal/log"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/version"
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

func CheckUpdate(logger *log.Logger, configuration *configuration.Configuration, version *version.Version) {
	lastCheckUnixTime := configuration.JSON.LastVersionCheck
	checkAllowedAfter := time.Now().Unix() - 3600 // 1h

	if lastCheckUnixTime > checkAllowedAfter {
		logger.Debug("CHECK UPDATE - Next check in %d seconds", lastCheckUnixTime-checkAllowedAfter)
		return
	}

	githubLastRelease, githubError := version.GetLatestRelease()
	if githubError != nil {
		logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		return
	}

	if version.GetCurrentVersion() != githubLastRelease.GetTagName() {
		logger.Info("⚠️  A new version of Gira is available! Check %s", "https://github.com/Ealenn/gira")
	}

	configuration.VersionChecked()
}
