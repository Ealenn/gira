package branch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type Manager struct {
	logger  *log.Logger
	git     *git.Git
	tracker issue.Tracker
}

func NewBranchManager(logger *log.Logger, git *git.Git, tracker issue.Tracker) *Manager {
	return &Manager{
		logger,
		git,
		tracker,
	}
}

func (manager *Manager) GetCurrentBranch() *Branch {
	currentBranch, currentBranchError := manager.git.CurrentBranch()
	if currentBranchError != nil {
		manager.logger.Fatal("❌ Unable to check current branch")
	}

	manager.logger.Debug("🔎 Current branch %s", currentBranch)
	branchNameParts := strings.Split(currentBranch, `/`)
	if len(branchNameParts) < 3 {
		manager.logger.Fatal("❌ Unable to find issue in branch name %s", currentBranch)
	}

	return &Branch{
		Type:    manager.getBranchType([]string{branchNameParts[0]}),
		IssueID: branchNameParts[1],
		Title:   branchNameParts[2],
	}
}

func (manager *Manager) Generate(issue *issue.Issue) *Branch {
	branchTitle := strings.ToLower(strings.TrimSpace(issue.Title))
	branchTitle = strings.Join(strings.Fields(branchTitle), "-")
	branchTitle = strings.ReplaceAll(branchTitle, " ", "-")
	branchTitle = regexp.MustCompile(`^\[[^\]]+\]\s*`).ReplaceAllString(branchTitle, "")
	branchTitle = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(branchTitle, "")
	branchTitle = regexp.MustCompile(`-+`).ReplaceAllString(branchTitle, "-")
	branchTitle = strings.Trim(branchTitle, "-")

	branchType := manager.getBranchType(issue.Types)

	branchRaw := fmt.Sprintf("%s/%s/%s", strings.ToLower(string(branchType)), strings.ToUpper(issue.ID), strings.ToLower(branchTitle))

	return &Branch{
		Type:    branchType,
		IssueID: issue.ID,
		Title:   branchTitle,
		Raw:     branchRaw,
	}
}

func (manager *Manager) getBranchType(issueTypes []string) Type {
	for _, issueType := range issueTypes {
		switch strings.ToLower(issueType) {
		case "bug":
			return Type(Bug)
		case "tasks":
		case "feature":
		case "enhancement":
			return Type(Feature)
		}
	}

	return Type(Feature)
}
