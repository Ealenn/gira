package cmd

import (
	Configuration "gira/configuration"
	"time"

	Services "gira/services"
)

func CmdVersion(configuration *Configuration.Configuration, loggerService *Services.LoggerService) {
	currentVersion := configuration.GetVersion(true)
	loggerService.Info("Current version : %s", currentVersion)

	githubService := Services.NewGitHubService(loggerService)
	githubLastRelease, githubError := githubService.GetLatestRelease()
	if githubError != nil {
		loggerService.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		loggerService.Fatal("❌ Failed to fetch latest release from GitHub, check %s", "https://github.com/Ealenn/gira")
		return
	}

	loggerService.Info("Latest version on GitHub: %s %s", githubLastRelease.TagName, Services.DebugStyle.Render("(", githubLastRelease.CreatedAt.Format(time.RFC822), ")"))

	if currentVersion == githubLastRelease.TagName {
		loggerService.Info("\n🚀 Gira is up to date.")
	} else {
		loggerService.Info("\n⚠️  A new version is available!\nCheck %s to update Gira", "https://github.com/Ealenn/gira")
	}
}
