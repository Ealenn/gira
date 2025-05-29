package ui

import (
	"time"

	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/service"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/version"
)

func CheckUpdate(logger *log.Logger, configuration *configuration.Configuration, version *version.Version) {
	lastCheckUnixTime := configuration.JSON.LastVersionCheck
	checkAllowedAfter := time.Now().Unix() - 3600 // 1h

	if lastCheckUnixTime > checkAllowedAfter {
		logger.Debug("CHECK UPDATE - Next check in %d seconds", lastCheckUnixTime-checkAllowedAfter)
		return
	}

	githubService := service.NewGitHub(logger)
	githubLastRelease, githubError := githubService.GetLatestRelease()
	if githubError != nil {
		logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		return
	}

	if version.GetCurrentVersion() != githubLastRelease.TagName {
		logger.Info("⚠️  A new version of Gira is available! Check %s", "https://github.com/Ealenn/gira")
	}

	configuration.VersionChecked()
}
