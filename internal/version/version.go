package version

import (
	"context"
	_ "embed"
	"strings"

	"github.com/Ealenn/gira/internal/log"
	"github.com/google/go-github/v73/github"
)

//go:embed version
var currentVersion string

type Version struct {
	logger       *log.Logger
	githubClient *github.Client
}

func New(logger *log.Logger) *Version {
	githubClient := github.NewClient(nil)

	return &Version{
		logger,
		githubClient,
	}
}

func (version *Version) GetCurrentVersion() string {
	return strings.TrimSpace(currentVersion)
}

func (version *Version) GetLatestRelease() (*github.RepositoryRelease, error) {
	release, releaseResponse, err := version.githubClient.Repositories.GetLatestRelease(context.Background(), "ealenn", "gira")

	if err != nil {
		version.logger.Debug("Unable to fetch Github release due to %s", releaseResponse.Status)
		return nil, err
	}

	return release, nil
}
