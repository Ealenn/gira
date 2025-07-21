package configuration

type ProfileType string

const (
	ProfileTypeJira   ProfileType = "JIRA"
	ProfileTypeGithub ProfileType = "GITHUB"
)

type JSONConfiguration struct {
	Profiles         []Profile `json:"profiles"`
	LastVersionCheck int64     `json:"lastVersionCheck,omitempty"`
}

type Profile struct {
	Name   string      `json:"name,omitempty"`
	Type   ProfileType `json:"type,omitempty"`
	Jira   Jira        `json:"jira,omitempty"`
	Github Github      `json:"github,omitempty"`
}

type Jira struct {
	Host      string `json:"host,omitempty"`
	Token     string `json:"token,omitempty"`
	AccountID string `json:"accountID,omitempty"`
	Email     string `json:"email,omitempty"`
	UserKey   string `json:"userKey,omitempty"`
}

type Github struct {
	User  string `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
