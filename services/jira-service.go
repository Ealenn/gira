package services

import (
	"context"
	"crypto/tls"
	Configuration "gira/configuration"
	"log"
	"net/http"

	Jira "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type JiraService struct {
	Client *Jira.Client
}

func NewJiraService(configuration *Configuration.Configuration) *JiraService {
	client, err := Jira.New(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, configuration.JSON.JiraHost)

	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	client.Auth.SetBearerToken(configuration.JSON.JiraToken)

	return &JiraService{
		Client: client,
	}
}

func (jiraService *JiraService) GetIssue(issueKeyId string) (*models.IssueSchemeV2, *models.ResponseScheme) {
	issue, response, err := jiraService.Client.Issue.Get(context.Background(), issueKeyId, nil, nil)

	if err != nil {
		return nil, response
	}

	return issue, response
}
