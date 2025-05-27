package commands

import (
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
	"github.com/Ealenn/gira/internal/services"
	"github.com/Ealenn/gira/internal/version"
)

type Version struct {
	logger        *logs.Logger
	configuration *configuration.Configuration
	version       *version.Version
	profile       *configuration.Profile
	githubService *services.GitHubService
}

func NewVersion(logger *logs.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Version {
	return &Version{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
		githubService: services.NewGitHubService(logger),
	}
}

func (command Version) Run() {
	currentVersion := command.version.GetCurrentVersion()
	command.logger.Info("Current version : %s", currentVersion)

	githubLastRelease, githubError := command.githubService.GetLatestRelease()
	if githubError != nil {
		command.logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		command.logger.Fatal("❌ Failed to fetch latest release from GitHub, check %s", "https://github.com/Ealenn/gira")
		return
	}

	command.logger.Info("Latest version on GitHub: %s %s", githubLastRelease.TagName, logs.DebugStyle.Render("(", githubLastRelease.CreatedAt.Format(time.RFC822), ")"))

	if currentVersion == githubLastRelease.TagName {
		command.logger.Info("\n🚀 Gira is up to date.")
	} else {
		command.logger.Info("\n⚠️  A new version is available!\nCheck %s to update Gira", "https://github.com/Ealenn/gira")
	}

	command.configuration.VersionChecked()
}
