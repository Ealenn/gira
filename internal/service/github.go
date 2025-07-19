package service

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Ealenn/gira/internal/log"

	"github.com/google/go-github/v73/github"
)

type GitHub struct {
	logger *log.Logger
	client *github.Client
}

func NewGitHub(logger *log.Logger) *GitHub {
	client := github.NewClient(nil)

	return &GitHub{
		logger,
		client,
	}
}

func (github *GitHub) GetLatestRelease() (*github.RepositoryRelease, error) {
	release, releaseResponse, err := github.client.Repositories.GetLatestRelease(context.Background(), "ealenn", "gira")

	if err != nil {
		github.logger.Debug("Unable to fetch Github release due to %s", releaseResponse.Status)
		return nil, err
	}

	return release, nil
}

func (github *GitHub) GetIssue(issueKeyID int) (*github.Issue, error) {
	username, repository := github.getCurrentRepository()
	issue, issueResponse, err := github.client.Issues.Get(context.Background(), username, repository, issueKeyID)

	if err != nil {
		github.logger.Debug("Unable to fetch Github release due to %s", issueResponse.Status)
		return nil, err
	}

	return issue, nil
}

func (github *GitHub) getCurrentRepository() (string, string) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		github.logger.Fatal("Error getting git remote URL")
	}

	repoURL := strings.TrimSpace(string(output))
	github.logger.Debug("Git repository URL:", repoURL)

	if !strings.HasPrefix(repoURL, "git@github.com:") || !strings.HasSuffix(repoURL, ".git") {
		github.logger.Fatal("Invalid Git Remote Origin format")
	}

	trimmed := strings.TrimPrefix(repoURL, "git@github.com:")
	trimmed = strings.TrimSuffix(trimmed, ".git")

	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		github.logger.Fatal("Invalid Git Remote Origin URL")
	}

	username := parts[0]
	repository := parts[1]

	return username, repository
}
