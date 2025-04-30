package cmd

import (
	Configuration "gira/configuration"
	UI "gira/ui"
	"regexp"
	"strings"

	Services "gira/services"
)

func CmdIssue(configuration *Configuration.Configuration, loggerService *Services.LoggerService, optionalIssueId *string) {
	UI.CheckConfiguration(loggerService, configuration)
	UI.CheckUpdate(loggerService, configuration)

	gitService := Services.NewGitService(loggerService)

	var issueId string

	if optionalIssueId != nil {
		issueId = *optionalIssueId
	} else {
		currentBranch, branchError := gitService.CurrentBranch()
		if branchError != nil {
			loggerService.Fatal("❌ Unable to check current branch")
		}
		loggerService.Debug("🔎 Current branch %s", currentBranch)

		branchNameParts := strings.Split(currentBranch, `/`)
		if len(branchNameParts) < 2 {
			loggerService.Fatal("❌ Unable to find issue in branch name %s", currentBranch)
		}

		issueId = branchNameParts[1]
	}

	loggerService.Debug("🔎 Issue %s", issueId)

	jiraService := Services.NewJiraService(configuration)
	issue, issueResponse := jiraService.GetIssue(issueId)
	if issue == nil {
		loggerService.Debug("Issue %s response status %s", issueId, issueResponse.Status)
		loggerService.Fatal("❌ Unable to find Jira %s", issueId)
	}

	loggerService.Log("%s: %s", Services.InfoStyle.Render("Summary"), issue.Fields.Summary)
	loggerService.Log("%s: %s - %s: %s - %s: %s", Services.InfoStyle.Render("Type"), issue.Fields.IssueType.Name, Services.InfoStyle.Render("Priority"), issue.Fields.Priority.Name, Services.InfoStyle.Render("Status"), issue.Fields.Status.Name)
	loggerService.Log("%s: %s <%s>", Services.InfoStyle.Render("Assignee"), issue.Fields.Assignee.DisplayName, issue.Fields.Assignee.EmailAddress)

	description := regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(string(issue.Fields.Description), "$1 $2")
	description = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(description, "$1")
	loggerService.Log("%s: \n%s", Services.InfoStyle.Render("Description"), description)

	loggerService.Info("\n🔗 More %s%s%s", configuration.JSON.JiraHost, "/browse/", issue.Key)
}
