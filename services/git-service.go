package services

import (
	"os/exec"
)

type GitService struct {
	loggerService *LoggerService
}

func NewGitService(loggerService *LoggerService) *GitService {
	return &GitService{
		loggerService,
	}
}

func (gitService *GitService) CreateBranch(name string) []byte {
	cmd := exec.Command("git", "checkout", "-b", name)
	output, err := cmd.CombinedOutput()

	if err != nil {
		gitService.loggerService.Fatal("%s", output)
	}

	return output
}
