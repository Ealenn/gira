package command

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/service"

	"github.com/manifoldco/promptui"
)

type Branch struct {
	logger        *log.Logger
	configuration *configuration.Configuration
	profile       *configuration.Profile
	jiraService   *service.Jira
	githubService *service.GitHub
	gitService    *service.Git
}

func NewBranch(logger *log.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Branch {
	return &Branch{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
		jiraService:   service.NewJira(logger, profile),
		githubService: service.NewGitHub(logger),
		gitService:    service.NewGit(logger),
	}
}

func (command Branch) Run(issueID string, assign bool, force bool) {
	var branchName string

	if command.profile.Type == configuration.ProfileTypeJira {
		branchName = command.getJiraBrancheName(issueID)
	} else {
		branchName = command.getGithubBrancheName(issueID)
	}

	if command.gitService.IsBranchExist(branchName) {
		command.logger.Warn("❌ Branch named %s already exists", branchName)

		if !force {
			command.logger.Info("Would you like to switch to this branch?\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		command.gitService.SwitchBranch(branchName)
		command.logger.Info("✅ %s has just been checkout", branchName)
	} else {
		if !force {
			command.logger.Info("🌳 Branch will be generated\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			prompt := promptui.Prompt{
				Label:     "Branch",
				Default:   branchName,
				Pointer:   promptui.PipeCursor,
				AllowEdit: true,
			}
			newBranchName, err := prompt.Run()
			if err != nil {
				command.logger.Fatal("The operation was %s", "canceled")
			}
			branchName = newBranchName
		}

		command.gitService.CreateBranch(branchName)
		command.logger.Info("✅ %s branch was created", branchName)
	}

	if assign {
		if assignError := command.jiraService.AssignIssue(issueID); assignError != nil {
			command.logger.Debug("%v", assignError)
			command.logger.Info("❌ %s Unable to assign issue %s to %s...", log.ErrorStyle.Render("Oups..."), issueID, command.profile.Jira.UserKey)
		} else {
			command.logger.Info("✅ Jira %s has been assigned to %s", issueID, command.profile.Jira.UserKey)
		}
	}
}

func (command Branch) getJiraBrancheName(issueID string) string {
	issue, issueResponse := command.jiraService.GetIssue(issueID)

	if issue == nil {
		command.logger.Debug("Issue %s response status %s", issueID, issueResponse.Status)
		command.logger.Fatal("❌ Unable to find issue %s", issueID)
	}

	return getBranchTitle(issue.Key, []string{issue.Fields.IssueType.Name}, issue.Fields.Summary)
}

func (command Branch) getGithubBrancheName(issueID string) string {
	issueNumber, err := strconv.Atoi(issueID)
	if err != nil {
		command.logger.Fatal("Github issue ID %s invalid !", issueID)
	}

	issue, issueErr := command.githubService.GetIssue(issueNumber)

	if issueErr != nil {
		command.logger.Fatal("❌ Unable to find issue %s", issueID)
	}

	labels := []string{}
	for _, element := range issue.Labels {
		labels = append(labels, *element.Name)
	}
	return getBranchTitle(issueID, labels, issue.GetTitle())
}

func getBranchType(issueTypeNames []string) string {
	for _, element := range issueTypeNames {
		switch strings.ToLower(element) {
		case "tasks":
		case "feature":
		case "enhancement":
			return "feature"
		case "bug":
			return "bugfix"
		}
	}

	return "feature"
}

func getBranchTitle(issueKeyID string, issueTypeNames []string, summary string) string {
	summary = strings.ToLower(strings.TrimSpace(summary))
	summary = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(summary, "")
	summary = strings.Join(strings.Fields(summary), " ")
	summary = strings.ReplaceAll(summary, " ", "-")
	return fmt.Sprintf("%s/%s/%s", getBranchType(issueTypeNames), issueKeyID, summary)
}
