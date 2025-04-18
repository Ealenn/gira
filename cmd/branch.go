package cmd

import (
	"bufio"
	"fmt"
	Configuration "gira/configuration"
	UI "gira/ui"
	"os"
	"regexp"
	"strings"

	Services "gira/services"
)

func CmdBranch(configuration *Configuration.Configuration, loggerService *Services.LoggerService, issueId string, assignIssue bool) {
	if !configuration.IsValid() {
		UI.ConfigurationError(configuration.Path)
	}

	jiraService := Services.NewJiraService(configuration)
	gitService := Services.NewGitService(loggerService)

	issue, issueResponse := jiraService.GetIssue(issueId)

	if issue == nil {
		loggerService.Debug("Issue %s response status %s", issueId, issueResponse.Status)
		loggerService.Fatal("Unable to find Jira %s", issueId)
	}

	branchName := getBranchTitle(issue.Key, issue.Fields.IssueType.Name, issue.Fields.Summary)

	loggerService.Info("Branch %s will be generated\nPress %s to continue, %s to cancel",
		UI.CodeStyle.Render(branchName),
		UI.InfoStyle.Render("ENTER"),
		UI.ErrorStyle.Render("CTRL+C"))
	bufio.NewReader(os.Stdin).ReadLine()

	gitService.CreateBranch(branchName)

	if assignIssue {
		if jiraUserAccount, jiraUserAccountError := jiraService.GetMyself(); jiraUserAccountError == nil {
			if assignError := jiraService.AssignIssueTo(issueId, jiraUserAccount.AccountID); assignError != nil {
				loggerService.Info("Unable to assign issue %s due to %v", issueId, assignError)
			} else {
				loggerService.Info("Jira %s has been assigned to %s", UI.InfoStyle.Render(issueId), UI.InfoStyle.Render(jiraUserAccount.DisplayName))
			}
		} else {
			loggerService.Info("Unable to fetch Jira account due to %v", jiraUserAccountError)
		}
	}
}

func getBranchType(issueTypeName string) string {
	switch issueTypeName {
	case "Tasks":
		return "feature"
	case "Bug":
		return "bugfix"
	default:
		return "feature"
	}
}

func getBranchTitle(issueKeyId string, issueTypeName string, summary string) string {
	summary = strings.ToLower(strings.TrimSpace(summary))
	summary = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(summary, "")
	summary = strings.Join(strings.Fields(summary), " ")
	summary = strings.ReplaceAll(summary, " ", "-")
	return fmt.Sprintf("%s/%s/%s", getBranchType(issueTypeName), issueKeyId, summary)
}
