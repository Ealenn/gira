package commands

import (
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
	"github.com/Ealenn/gira/internal/services"
	"github.com/Ealenn/gira/internal/ui"
)

func Issue(configuration *configuration.Configuration, logger *logs.Logger, optionalIssueID *string) {
	ui.CheckConfiguration(logger, configuration)
	ui.CheckUpdate(logger, configuration)

	gitService := services.NewGitService(logger)

	var issueID string

	if optionalIssueID != nil {
		issueID = *optionalIssueID
	} else {
		currentBranch, branchError := gitService.CurrentBranch()
		if branchError != nil {
			logger.Fatal("❌ Unable to check current branch")
		}
		logger.Debug("🔎 Current branch %s", currentBranch)

		branchNameParts := strings.Split(currentBranch, `/`)
		if len(branchNameParts) < 2 {
			logger.Fatal("❌ Unable to find issue in branch name %s", currentBranch)
		}

		issueID = branchNameParts[1]
	}

	logger.Debug("🔎 Issue %s", issueID)

	jiraService := services.NewJiraService(configuration)
	issue, issueResponse := jiraService.GetIssue(issueID)
	if issue == nil {
		logger.Debug("Issue %s response status %s", issueID, issueResponse.Status)
		logger.Fatal("❌ Unable to find Jira %s", issueID)
	}

	logger.Log("%s: %s", logs.InfoStyle.Render("Summary"), issue.Fields.Summary)
	logger.Log("%s: %s - %s: %s - %s: %s", logs.InfoStyle.Render("Type"), issue.Fields.IssueType.Name, logs.InfoStyle.Render("Priority"), issue.Fields.Priority.Name, logs.InfoStyle.Render("Status"), issue.Fields.Status.Name)
	logger.Log("%s: %s <%s>", logs.InfoStyle.Render("Assignee"), issue.Fields.Assignee.DisplayName, issue.Fields.Assignee.EmailAddress)

	description := regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(string(issue.Fields.Description), "$1 $2")
	description = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(description, "$1")
	logger.Log("%s: \n%s", logs.InfoStyle.Render("Description"), description)

	logger.Info("\n🔗 More %s%s%s", configuration.JSON.JiraHost, "/browse/", issue.Key)
}
