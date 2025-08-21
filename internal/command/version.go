package command

import (
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/version"
)

type Version struct {
	logger        *log.Logger
	configuration *configuration.Configuration
	version       *version.Version
}

func NewVersion(logger *log.Logger, configuration *configuration.Configuration, version *version.Version) *Version {
	return &Version{
		logger,
		configuration,
		version,
	}
}

func (cmd Version) Run() {
	currentVersion := cmd.version.GetCurrentVersion()
	cmd.logger.Info("Current version : %s", currentVersion)

	githubLastRelease, githubError := cmd.version.GetLatestRelease()
	if githubError != nil {
		cmd.logger.Debug("Unable to fetch latest version of Gira due to %v", githubError)
		cmd.logger.Fatal("❌ Failed to fetch latest release from GitHub, check %s", "https://github.com/Ealenn/gira")
		return
	}

	latestTag := githubLastRelease.GetTagName()
	cmd.logger.Info("Latest version on GitHub: %s %s", latestTag, log.DebugStyle.Render("(", githubLastRelease.CreatedAt.Format(time.RFC822), ")"))

	if currentVersion == latestTag {
		cmd.logger.Info("\n🚀 Gira is up to date.")
	} else {
		cmd.logger.Info("\n⚠️  A new version is available!\nCheck %s to update Gira", "https://github.com/Ealenn/gira")
	}

	cmd.configuration.VersionChecked()
}
