package ui

import (
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
	"github.com/Ealenn/gira/internal/services"
)

func CheckUpdate(logger *logs.Logger, configuration *configuration.Configuration) {
	lastCheckUnixTime := configuration.JSON.GiraLastVersionCheck
	checkAllowedAfter := time.Now().Unix() - 3600 // 1h

	if lastCheckUnixTime > checkAllowedAfter {
		logger.Debug("UPDATE - Last check : %d", lastCheckUnixTime)
		logger.Debug("UPDATE - Check allowed after : %d", checkAllowedAfter)
		logger.Debug("UPDATE - Next check in %d seconds", lastCheckUnixTime-checkAllowedAfter)
		return
	}

	currentVersion := configuration.GetVersion(true)

	githubService := services.NewGitHubService(logger)
	githubLastRelease, githubError := githubService.GetLatestRelease()
	if githubError != nil {
		logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		return
	}

	if currentVersion != githubLastRelease.TagName {
		logger.Info("⚠️  A new version of Gira is available! Check %s", "https://github.com/Ealenn/gira")
	}
}
