package command

import (
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/browser"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type Open struct {
	logger  *log.Logger
	branch  *branch.Manager
	tracker issue.Tracker
}

func NewOpen(logger *log.Logger, branch *branch.Manager, tracker issue.Tracker) *Open {
	return &Open{
		logger,
		branch,
		tracker,
	}
}

func (command Open) Run(optionalIssueID *string) {
	browser := browser.NewBrowser(command.logger)

	var issueID string
	if optionalIssueID != nil {
		issueID = *optionalIssueID
	} else {
		issueID = command.branch.GetCurrentBranch().IssueID
	}

	issue := command.tracker.GetIssue(issueID)
	command.logger.Info("Open Issue %s : %s", issue.ID, issue.Title)
	browser.Open(issue.URL)
}
