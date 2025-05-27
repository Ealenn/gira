package version

import (
	_ "embed"
	"strings"
)

//go:embed version
var currentVersion string

type Version struct{}

func New() *Version {
	return &Version{}
}

func (version *Version) GetCurrentVersion() string {
	return strings.TrimSpace(currentVersion)
}
