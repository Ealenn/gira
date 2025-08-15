package issue

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
}

type CreateIssueOptions struct {
	Title       string
	Description string
	Type        Type
	Project     string
}

type Tracker interface {
	GetIssue(issueKeyID string) *Issue
	CreateIssue(options CreateIssueOptions) *Issue
	SelfAssignIssue(issueKeyID string) error
}
