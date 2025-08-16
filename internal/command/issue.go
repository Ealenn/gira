package command

import (
	"regexp"

	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type Issue struct {
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
}

func NewIssue(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Issue {
	return &Issue{
		logger,
		tracker,
		git,
		branch,
	}
}

func (command Issue) Run(optionalIssueID *string, enableAI bool) {
	var issueID string

	if optionalIssueID != nil {
		issueID = *optionalIssueID
	} else {
		issueID = command.branch.GetCurrentBranch().IssueID
	}

	issue := command.tracker.GetIssue(issueID)

	command.logger.Log("%s: %s", log.InfoStyle.Render("Summary"), issue.Title)
	command.logger.Log("%s: %s - %s: %s", log.InfoStyle.Render("Type"), issue.Types, log.InfoStyle.Render("State"), issue.Status)
	for _, assignee := range issue.Assignees {
		command.logger.Log("%s: %s <%s>", log.InfoStyle.Render("Assignee"), assignee.Name, assignee.Email)
	}

	description := regexp.MustCompile(`\[(.*?)\|(.*?)\]`).ReplaceAllString(issue.Description, "$1 $2")
	description = regexp.MustCompile(`\[(.*?)\]`).ReplaceAllString(description, "$1")
	command.logger.Log("%s: \n%s", log.InfoStyle.Render("Description"), description)

	if enableAI {
		agent := ai.NewOpenAI(command.logger)
		response, err := agent.IssueSummary(issue)

		if err == nil {
			command.logger.Log("--------------------------------------------------------------")
			command.logger.Info("%s", "AI Quick summary:")
			command.logger.Log("%s", response)
			command.logger.Log("--------------------------------------------------------------")
		} else {
			command.logger.Debug("Unable to generate summary due to %v", err)
		}
	}

	command.logger.Info("\n🔗 More %s", issue.URL)
}
