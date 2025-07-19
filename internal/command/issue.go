package command

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/service"
)

type Issue struct {
	logger        *log.Logger
	configuration *configuration.Configuration
	profile       *configuration.Profile
	gitService    *service.Git
	jiraService   *service.Jira
	githubService *service.GitHub
}

func NewIssue(logger *log.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Issue {
	return &Issue{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
		gitService:    service.NewGit(logger),
		jiraService:   service.NewJira(logger, profile),
		githubService: service.NewGitHub(logger),
	}
}

func (command Issue) Run(optionalIssueID *string) {
	var issueID string

	if optionalIssueID != nil {
		issueID = *optionalIssueID
	} else {
		currentBranch, branchError := command.gitService.CurrentBranch()
		if branchError != nil {
			command.logger.Fatal("❌ Unable to check current branch")
		}
		command.logger.Debug("🔎 Current branch %s", currentBranch)

		branchNameParts := strings.Split(currentBranch, `/`)
		if len(branchNameParts) < 2 {
			command.logger.Fatal("❌ Unable to find issue in branch name %s", currentBranch)
		}

		issueID = branchNameParts[1]
		command.logger.Debug("🔎 Issue %s", issueID)
	}

	if command.profile.Type == configuration.ProfileTypeJira {
		command.showJiraIssue(issueID)
	}

	if command.profile.Type == configuration.ProfileTypeGithub {
		command.showGithubIssue(issueID)
	}
}

func (command Issue) showJiraIssue(issueID string) {
	issue, issueResponse := command.jiraService.GetIssue(issueID)
	if issue == nil {
		command.logger.Debug("Issue %s response status %s", issueID, issueResponse.Status)
		command.logger.Fatal("❌ Unable to find issue %s", issueID)
	}

	command.logger.Log("%s: %s", log.InfoStyle.Render("Summary"), issue.Fields.Summary)
	command.logger.Log("%s: %s - %s: %s - %s: %s", log.InfoStyle.Render("Type"), issue.Fields.IssueType.Name, log.InfoStyle.Render("Priority"), issue.Fields.Priority.Name, log.InfoStyle.Render("Status"), issue.Fields.Status.Name)
	if issue.Fields.Assignee != nil {
		command.logger.Log("%s: %s <%s>", log.InfoStyle.Render("Assignee"), issue.Fields.Assignee.DisplayName, issue.Fields.Assignee.EmailAddress)
	}

	description := regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(string(issue.Fields.Description), "$1 $2")
	description = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(description, "$1")
	command.logger.Log("%s: \n%s", log.InfoStyle.Render("Description"), description)

	command.logger.Info("\n🔗 More %s%s%s", command.profile.Jira.Host, "/browse/", issue.Key)
}

func (command Issue) showGithubIssue(issueID string) {
	issueNumber, err := strconv.Atoi(issueID)
	if err != nil {
		command.logger.Fatal("Github issue ID %s invalid !", issueID)
	}

	issue, issueErr := command.githubService.GetIssue(issueNumber)
	if issueErr != nil {
		command.logger.Debug("Issue %s error %s", issueID, issueErr)
		command.logger.Fatal("❌ Unable to find issue %s", issueID)
	}

	command.logger.Log("%s: %s", log.InfoStyle.Render("Summary"), issue.GetTitle())

	labels := []string{}
	for _, element := range issue.Labels {
		labels = append(labels, *element.Name)
	}
	command.logger.Log("%s: %s - %s: %s", log.InfoStyle.Render("Type"), labels, log.InfoStyle.Render("State"), issue.GetState())
	if issue.GetAssignee() != nil {
		command.logger.Log("%s: %s <%s>", log.InfoStyle.Render("Assignee"), issue.GetAssignee().GetLogin(), issue.GetAssignee().GetEmail())
	}

	description := regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(string(issue.GetBody()), "$1 $2")
	description = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(description, "$1")
	command.logger.Log("%s: \n%s", log.InfoStyle.Render("Description"), description)

	command.logger.Info("\n🔗 More %s", issue.GetHTMLURL())
}
