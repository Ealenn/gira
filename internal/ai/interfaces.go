package ai

import "github.com/Ealenn/gira/internal/issue"

type Agent interface {
	BranchNames(issue *issue.Issue) ([]string, error)
	CommitNames(issue *issue.Issue) ([]string, error)
	IssueSummary(issue *issue.Issue) (string, error)
	Rewrite(context string, text string) (string, error)
}
