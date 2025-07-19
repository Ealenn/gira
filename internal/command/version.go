package command

import (
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/service"
	"github.com/Ealenn/gira/internal/version"
)

type Version struct {
	logger        *log.Logger
	configuration *configuration.Configuration
	version       *version.Version
	profile       *configuration.Profile
	githubService *service.GitHub
}

func NewVersion(logger *log.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Version {
	return &Version{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
		githubService: service.NewGitHub(logger),
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

	latestTag := githubLastRelease.GetTagName()
	command.logger.Info("Latest version on GitHub: %s %s", latestTag, log.DebugStyle.Render("(", githubLastRelease.CreatedAt.Format(time.RFC822), ")"))

	if currentVersion == latestTag {
		command.logger.Info("\n🚀 Gira is up to date.")
	} else {
		command.logger.Info("\n⚠️  A new version is available!\nCheck %s to update Gira", "https://github.com/Ealenn/gira")
	}

	command.configuration.VersionChecked()
}
