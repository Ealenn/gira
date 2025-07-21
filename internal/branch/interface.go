package branch

type BranchType string

const (
	Bug     BranchType = "BUG"
	Feature BranchType = "FEATURE"
)

type Branch struct {
	Type    BranchType
	IssueID string
	Title   string
	Raw     string
}
