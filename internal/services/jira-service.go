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
	Client        *v2.Client
	Configuration *configuration.Configuration
}

func NewJiraService(configuration *configuration.Configuration) *JiraService {
	client, err := v2.New(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, configuration.JSON.JiraHost)

	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	if configuration.JSON.JiraEmail != "" {
		client.Auth.SetBasicAuth(configuration.JSON.JiraEmail, configuration.JSON.JiraToken)
	} else {
		client.Auth.SetBearerToken(configuration.JSON.JiraToken)
	}

	return &JiraService{
		Client:        client,
		Configuration: configuration,
	}
}

func (jiraService *JiraService) GetIssue(issueKeyID string) (*models.IssueSchemeV2, *models.ResponseScheme) {
	issue, response, err := jiraService.Client.Issue.Get(context.Background(), issueKeyID, nil, nil)

	if err != nil {
		return nil, response
	}

	return issue, response
}

func (jiraService *JiraService) GetMyself() (*models.UserScheme, error) {
	user, _, userError := jiraService.Client.MySelf.Details(context.Background(), []string{})
	if userError != nil {
		return nil, userError
	}
	return user, nil
}

func (jiraService *JiraService) AssignIssue(issueKeyID string) error {
	ctx := context.Background()

	// For Jira Cloud, use AccountID
	if jiraService.Configuration.JSON.JiraAccountID != "" {
		_, err := jiraService.Client.Issue.Assign(ctx, issueKeyID, jiraService.Configuration.JSON.JiraAccountID)
		return err
	}

	// For Jira Server/Data Center, use User Key or Name
	_, err := jiraService.Client.Issue.Update(ctx, issueKeyID, true, &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Assignee: &models.UserScheme{
				Key:  jiraService.Configuration.JSON.JiraUserKey,
				Name: jiraService.Configuration.JSON.JiraUserKey,
			},
		},
	}, nil, nil)

	return err
}
