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
)

type Branch struct {
	logger        *logs.Logger
	configuration *configuration.Configuration
	profile       *configuration.Profile
	jiraService   *services.JiraService
	gitService    *services.GitService
}

func NewBranch(logger *logs.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Branch {
	return &Branch{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
		jiraService:   services.NewJiraService(profile),
		gitService:    services.NewGitService(logger),
	}
}

func (command Branch) Run(issueID string, assign bool, force bool) {
	issue, issueResponse := command.jiraService.GetIssue(issueID)

	if issue == nil {
		command.logger.Debug("Issue %s response status %s", issueID, issueResponse.Status)
		command.logger.Fatal("❌ Unable to find issue %s", issueID)
	}

	branchName := getBranchTitle(issue.Key, issue.Fields.IssueType.Name, issue.Fields.Summary)
	if command.gitService.IsBranchExist(branchName) {
		command.logger.Warn("❌ Branch named %s already exists", branchName)

		if !force {
			command.logger.Info("Would you like to switch to this branch?\nPress %s to continue, %s to cancel", "ENTER", logs.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		command.gitService.SwitchBranch(branchName)
		command.logger.Info("✅ %s has just been checkout", branchName)
	} else {
		if !force {
			command.logger.Info("🌳 Branch %s will be generated\nPress %s to continue, %s to cancel", branchName, "ENTER", logs.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		command.gitService.CreateBranch(branchName)
		command.logger.Info("✅ %s branch was created", branchName)
	}

	if assign {
		if assignError := command.jiraService.AssignIssue(issueID); assignError != nil {
			command.logger.Debug("%v", assignError)
			command.logger.Info("❌ %s Unable to assign issue %s to %s...", logs.ErrorStyle.Render("Oups..."), issueID, command.profile.Jira.UserKey)
		} else {
			command.logger.Info("✅ Jira %s has been assigned to %s", issueID, command.profile.Jira.UserKey)
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
