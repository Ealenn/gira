package issue

import (
	"context"
	"strconv"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/log"

	"github.com/google/go-github/v73/github"
)

type GitHubTracker struct {
	logger       *log.Logger
	profile      *configuration.Profile
	git          *git.Git
	githubClient *github.Client
}

func NewGitHub(logger *log.Logger, profile *configuration.Profile, git *git.Git) *GitHubTracker {
	client := github.NewClient(nil)

	if profile.Github.Token != "" {
		client = client.WithAuthToken(profile.Github.Token)
	}

	return &GitHubTracker{
		logger,
		profile,
		git,
		client,
	}
}

func (tracker *GitHubTracker) GetIssue(issueKeyID string) *Issue {
	issueNumber := tracker.getIssueNumber(issueKeyID)
	username, repository := tracker.getCurrentRepository()
	issue, issueResponse, err := tracker.githubClient.Issues.Get(context.Background(), username, repository, issueNumber)

	if err != nil {
		tracker.logger.Debug("Issue %s response status %s with error %v", issueKeyID, issueResponse.Status, err)
		tracker.logger.Fatal("❌ Unable to find issue %s", issueKeyID)
	}

	return tracker.formatIssue(issue)
}

func (tracker *GitHubTracker) CreateIssue(options CreateIssueOptions) *Issue {
	username, repository := tracker.getCurrentRepository()

	issue, response, err := tracker.githubClient.Issues.Create(context.Background(), username, repository, &github.IssueRequest{
		Title: &options.Title,
		Body:  &options.Description,
	})

	if err != nil {
		tracker.logger.Fatal("Unable to create Github issue due to %s", response.Status)
	}

	tracker.logger.Debug("Issue %v created", issue.GetNumber())
	return tracker.formatIssue(issue)
}

func (tracker *GitHubTracker) SelfAssignIssue(issueKeyID string) error {
	issueNumber := tracker.getIssueNumber(issueKeyID)
	username, repository := tracker.getCurrentRepository()
	_, issueResponse, err := tracker.githubClient.Issues.AddAssignees(context.Background(), username, repository, issueNumber, []string{tracker.profile.Github.User})

	if err != nil {
		tracker.logger.Debug("Unable to fetch Github release due to %s", issueResponse.Status)
		return err
	}

	return nil
}

func (tracker *GitHubTracker) getCurrentRepository() (string, string) {
	origin := tracker.git.CurrentOrigin()

	if !strings.HasPrefix(origin, "git@github.com:") || !strings.HasSuffix(origin, ".git") {
		tracker.logger.Fatal("Invalid Git Remote Origin format")
	}

	trimmed := strings.TrimPrefix(origin, "git@github.com:")
	trimmed = strings.TrimSuffix(trimmed, ".git")

	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		tracker.logger.Fatal("Invalid Git Remote Origin URL")
	}

	username := parts[0]
	repository := parts[1]

	return username, repository
}

func (tracker *GitHubTracker) getIssueNumber(issueKeyID string) int {
	issueNumber, err := strconv.Atoi(issueKeyID)
	if err != nil {
		tracker.logger.Fatal("Github issue ID %s invalid !", issueKeyID)
	}
	return issueNumber
}

func (tracker *GitHubTracker) getIssueString(issueKeyID int) string {
	return strconv.Itoa(issueKeyID)
}

func (tracker *GitHubTracker) formatIssue(issue *github.Issue) *Issue {
	var assignees []Assignee
	if issue.Assignees != nil {
		for _, assignee := range issue.Assignees {
			email := assignee.GetEmail()
			if email == "" {
				email = assignee.GetHTMLURL()
			}
			assignees = append(assignees, Assignee{
				ID:    assignee.GetLogin(),
				Name:  assignee.GetLogin(),
				Email: email,
			})
		}
	}

	labels := []string{}
	for _, element := range issue.Labels {
		labels = append(labels, *element.Name)
	}

	return &Issue{
		ID:          tracker.getIssueString(issue.GetNumber()),
		Title:       issue.GetTitle(),
		Description: issue.GetBody(),
		Status:      issue.GetState(),
		Types:       labels,
		Assignees:   assignees,
		URL:         issue.GetHTMLURL(),
		CreatedAt:   issue.CreatedAt.Time,
	}
}
