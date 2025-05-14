package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/logs"
	"github.com/Ealenn/gira/internal/services"
	"github.com/Ealenn/gira/internal/ui"
)

func Branch(configuration *configuration.Configuration, logger *logs.Logger, issueID string, assignIssue bool, force bool) {
	ui.CheckConfiguration(logger, configuration)
	ui.CheckUpdate(logger, configuration)

	jiraService := services.NewJiraService(configuration)
	gitService := services.NewGitService(logger)

	issue, issueResponse := jiraService.GetIssue(issueID)

	if issue == nil {
		logger.Debug("Issue %s response status %s", issueID, issueResponse.Status)
		logger.Fatal("❌ Unable to find Jira %s", issueID)
	}

	branchName := getBranchTitle(issue.Key, issue.Fields.IssueType.Name, issue.Fields.Summary)
	if gitService.IsBranchExist(branchName) {
		logger.Warn("❌ Branch named %s already exists", branchName)

		if !force {
			logger.Info("Would you like to switch to this branch?\nPress %s to continue, %s to cancel", "ENTER", logs.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		gitService.SwitchBranch(branchName)
		logger.Info("✅ %s has just been checkout", branchName)
	} else {
		if !force {
			logger.Info("🌳 Branch %s will be generated\nPress %s to continue, %s to cancel", branchName, "ENTER", logs.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		gitService.CreateBranch(branchName)
		logger.Info("✅ %s branch was created", branchName)
	}

	if assignIssue {
		if assignError := jiraService.AssignIssue(issueID); assignError != nil {
			logger.Debug("%v", assignError)
			logger.Info("❌ %s Unable to assign issue %s to %s...", logs.ErrorStyle.Render("Oups..."), issueID, configuration.JSON.JiraUserKey)
		} else {
			logger.Info("✅ Jira %s has been assigned to %s", issueID, configuration.JSON.JiraUserKey)
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

func getBranchTitle(issueKeyID string, issueTypeName string, summary string) string {
	summary = strings.ToLower(strings.TrimSpace(summary))
	summary = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(summary, "")
	summary = strings.Join(strings.Fields(summary), " ")
	summary = strings.ReplaceAll(summary, " ", "-")
	return fmt.Sprintf("%s/%s/%s", getBranchType(issueTypeName), issueKeyID, summary)
}
