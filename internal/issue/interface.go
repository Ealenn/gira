package issue

type Type string

const (
	Bug     Type = "BUG"
	Feature Type = "FEATURE"
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
}

type Tracker interface {
	GetIssue(issueKeyID string) *Issue
	SelfAssignIssue(issueKeyID string) error
}
