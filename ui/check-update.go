package UI

import (
	"gira/configuration"
	"gira/services"
	"time"
)

func CheckUpdate(loggerService *services.LoggerService, configuration *configuration.Configuration) {
	lastCheckUnixTime := configuration.JSON.GiraLastVersionCheck
	checkAllowedAfter := time.Now().Unix() - 3600 // 1h

	if lastCheckUnixTime > checkAllowedAfter {
		loggerService.Debug("UPDATE - Last check : %d", lastCheckUnixTime)
		loggerService.Debug("UPDATE - Check allowed after : %d", checkAllowedAfter)
		loggerService.Debug("UPDATE - Next check in %d seconds", lastCheckUnixTime-checkAllowedAfter)
		return
	}

	currentVersion := configuration.GetVersion(true)

	githubService := services.NewGitHubService(loggerService)
	githubLastRelease, githubError := githubService.GetLatestRelease()
	if githubError != nil {
		loggerService.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		return
	}

	if currentVersion != githubLastRelease.TagName {
		loggerService.Info("⚠️  A new version of Gira is available! Check %s", "https://github.com/Ealenn/gira")
	}
}
