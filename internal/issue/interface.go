package issue

type IssueType string

const (
	Bug     IssueType = "BUG"
	Feature IssueType = "FEATURE"
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
