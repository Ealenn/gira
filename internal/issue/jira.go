package issue

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/log"

	v2 "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type JiraTracker struct {
	logger     *log.Logger
	profile    *configuration.Profile
	git        *git.Git
	jiraClient *v2.Client
}

func NewJira(logger *log.Logger, profile *configuration.Profile, git *git.Git) *JiraTracker {
	client, err := v2.New(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, profile.Jira.Host)

	if err != nil {
		logger.Debug("Jira client error: %v", err)
		logger.Fatal("Unable to create Jira Client")
	}

	client.Auth.SetBearerToken(profile.Jira.Token)

	return &JiraTracker{
		logger:     logger,
		profile:    profile,
		git:        git,
		jiraClient: client,
	}
}

func (tracker *JiraTracker) GetIssue(issueKeyID string) *Issue {
	issue, issueResponse, err := tracker.jiraClient.Issue.Get(context.Background(), issueKeyID, nil, nil)

	if err != nil {
		tracker.logger.Debug("Issue %s response status %s", issueKeyID, issueResponse.Status)
		tracker.logger.Fatal("❌ Unable to find issue %s", issueKeyID)
	}

	return tracker.formatIssue(issue)
}

func (tracker *JiraTracker) CreateIssue(options CreateIssueOptions) *Issue {
	issueTypeName := "Task"
	if options.Type == TypeBug {
		issueTypeName = "Bug"
	}

	issue, issueResponse, err := tracker.jiraClient.Issue.Create(context.Background(), &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			IssueType:   &models.IssueTypeScheme{Name: issueTypeName},
			Summary:     options.Title,
			Description: options.Description,
			Project:     &models.ProjectScheme{Key: options.Project},
		},
	}, &models.CustomFields{})

	if err != nil {
		tracker.logger.Debug("Create issue response status %s with error %v", issueResponse.StatusCode, err)
		tracker.logger.Fatal("❌ Unable to create issue due to %s", issueResponse.Status)
	}

	tracker.logger.Debug("Issue ID %s and Key %s created", issue.ID, issue.Key)
	return tracker.GetIssue(issue.Key)
}

func (tracker *JiraTracker) GetMyself() (*models.UserScheme, error) {
	user, _, userError := tracker.jiraClient.MySelf.Details(context.Background(), []string{})
	if userError != nil {
		return nil, userError
	}
	return user, nil
}

func (tracker *JiraTracker) SelfAssignIssue(issueKeyID string) error {
	ctx := context.Background()

	user, getMyselfErr := tracker.GetMyself()
	if getMyselfErr != nil {
		return getMyselfErr
	}

	// For Jira Cloud, use AccountID
	if user.AccountID != "" {
		_, err := tracker.jiraClient.Issue.Assign(ctx, issueKeyID, user.AccountID)
		return err
	}

	// For Jira Server/Data Center, use User Key or Name
	_, err := tracker.jiraClient.Issue.Update(ctx, issueKeyID, true, &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Assignee: &models.UserScheme{
				Key:  user.Key,
				Name: user.Key,
			},
		},
	}, nil, nil)

	return err
}

func (tracker *JiraTracker) formatIssue(issue *models.IssueSchemeV2) *Issue {
	var assignees []Assignee
	if issue.Fields.Assignee != nil {
		assignees = append(assignees, Assignee{
			ID:    issue.Fields.Assignee.AccountID,
			Name:  issue.Fields.Assignee.DisplayName,
			Email: issue.Fields.Assignee.EmailAddress,
		})
	}

	return &Issue{
		ID:          issue.Key,
		Title:       issue.Fields.Summary,
		Description: tracker.toMarkdown(issue.Fields.Description),
		Status:      issue.Fields.Status.Name,
		Types:       []string{issue.Fields.IssueType.Name},
		URL:         fmt.Sprintf("%s%s%s", tracker.profile.Jira.Host, "/browse/", issue.Key),
		Assignees:   assignees,
		CreatedAt:   time.Time(*issue.Fields.Created),
	}
}

func (tracker *JiraTracker) toMarkdown(content string) string {
	content = regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(content, "$1 $2")
	content = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(content, "$1")
	return content
}
