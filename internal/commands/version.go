package commands

import (
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
	"github.com/Ealenn/gira/internal/services"
)

func Version(configuration *configuration.Configuration, logger *logs.Logger) {
	currentVersion := configuration.GetVersion(true)
	logger.Info("Current version : %s", currentVersion)

	githubService := services.NewGitHubService(logger)
	githubLastRelease, githubError := githubService.GetLatestRelease()
	if githubError != nil {
		logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		logger.Fatal("❌ Failed to fetch latest release from GitHub, check %s", "https://github.com/Ealenn/gira")
		return
	}

	logger.Info("Latest version on GitHub: %s %s", githubLastRelease.TagName, logs.DebugStyle.Render("(", githubLastRelease.CreatedAt.Format(time.RFC822), ")"))

	if currentVersion == githubLastRelease.TagName {
		logger.Info("\n🚀 Gira is up to date.")
	} else {
		logger.Info("\n⚠️  A new version is available!\nCheck %s to update Gira", "https://github.com/Ealenn/gira")
	}
}
