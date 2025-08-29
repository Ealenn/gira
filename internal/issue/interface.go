package issue

import "time"

type Type string

const (
	TypeBug     Type = "BUG"
	TypeFeature Type = "FEATURE"
)

type Assignee struct {
	ID    string
	Name  string
	Email string
}

type Issue struct {
	ID          string
	Title       string
	Description string
	Status      string
	Types       []string
	Assignees   []Assignee
	URL         string
	CreatedAt   time.Time
}

type CreateIssueOptions struct {
	Title       string
	Description string
	Type        Type
	Project     string
}

type Tracker interface {
	SearchIssues(status string) map[string]*Issue
	GetIssue(issueKeyID string) *Issue
	CreateIssue(options CreateIssueOptions) *Issue
	SelfAssignIssue(issueKeyID string) error
}
