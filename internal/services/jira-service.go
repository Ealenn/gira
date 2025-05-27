package services

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Ealenn/gira/internal/configuration"

	v2 "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type JiraService struct {
	client  *v2.Client
	profile *configuration.Profile
}

func NewJiraService(profile *configuration.Profile) *JiraService {
	client, err := v2.New(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, profile.Jira.Host)

	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	if profile.Jira.Email != "" {
		client.Auth.SetBasicAuth(profile.Jira.Email, profile.Jira.Token)
	} else {
		client.Auth.SetBearerToken(profile.Jira.Token)
	}

	return &JiraService{
		client:  client,
		profile: profile,
	}
}

func (jiraService *JiraService) GetIssue(issueKeyID string) (*models.IssueSchemeV2, *models.ResponseScheme) {
	issue, response, err := jiraService.client.Issue.Get(context.Background(), issueKeyID, nil, nil)

	if err != nil {
		return nil, response
	}

	return issue, response
}

func (jiraService *JiraService) GetMyself() (*models.UserScheme, error) {
	user, _, userError := jiraService.client.MySelf.Details(context.Background(), []string{})
	if userError != nil {
		return nil, userError
	}
	return user, nil
}

func (jiraService *JiraService) AssignIssue(issueKeyID string) error {
	ctx := context.Background()

	// For Jira Cloud, use AccountID
	if jiraService.profile.Jira.AccountID != "" {
		_, err := jiraService.client.Issue.Assign(ctx, issueKeyID, jiraService.profile.Jira.AccountID)
		return err
	}

	// For Jira Server/Data Center, use User Key or Name
	_, err := jiraService.client.Issue.Update(ctx, issueKeyID, true, &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Assignee: &models.UserScheme{
				Key:  jiraService.profile.Jira.UserKey,
				Name: jiraService.profile.Jira.UserKey,
			},
		},
	}, nil, nil)

	return err
}
