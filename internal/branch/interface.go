package branch

type Type string

const (
	Bug     Type = "BUGFIX"
	Feature Type = "FEATURE"
)

type Branch struct {
	Type    Type
	IssueID string
	Title   string
	Raw     string
}
