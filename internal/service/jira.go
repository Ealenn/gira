package service

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"

	v2 "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

type Jira struct {
	logger  *log.Logger
	client  *v2.Client
	profile *configuration.Profile
}

func NewJira(logger *log.Logger, profile *configuration.Profile) *Jira {
	if profile.Type != configuration.ProfileTypeJira {
		return &Jira{
			logger:  logger,
			client:  nil,
			profile: profile,
		}
	}

	client, err := v2.New(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}, profile.Jira.Host)

	if err != nil {
		logger.Debug("Jira client error: %v", err)
		logger.Fatal("Unable to create Jira Client")
	}

	if profile.Jira.Email != "" {
		client.Auth.SetBasicAuth(profile.Jira.Email, profile.Jira.Token)
	} else {
		client.Auth.SetBearerToken(profile.Jira.Token)
	}

	return &Jira{
		logger:  logger,
		client:  client,
		profile: profile,
	}
}

func (jira *Jira) GetIssue(issueKeyID string) (*models.IssueSchemeV2, *models.ResponseScheme) {
	issue, response, err := jira.client.Issue.Get(context.Background(), issueKeyID, nil, nil)

	if err != nil {
		return nil, response
	}

	return issue, response
}

func (jira *Jira) GetMyself() (*models.UserScheme, error) {
	user, _, userError := jira.client.MySelf.Details(context.Background(), []string{})
	if userError != nil {
		return nil, userError
	}
	return user, nil
}

func (jira *Jira) AssignIssue(issueKeyID string) error {
	ctx := context.Background()

	// For Jira Cloud, use AccountID
	if jira.profile.Jira.AccountID != "" {
		_, err := jira.client.Issue.Assign(ctx, issueKeyID, jira.profile.Jira.AccountID)
		return err
	}

	// For Jira Server/Data Center, use User Key or Name
	_, err := jira.client.Issue.Update(ctx, issueKeyID, true, &models.IssueSchemeV2{
		Fields: &models.IssueFieldsSchemeV2{
			Assignee: &models.UserScheme{
				Key:  jira.profile.Jira.UserKey,
				Name: jira.profile.Jira.UserKey,
			},
		},
	}, nil, nil)

	return err
}
