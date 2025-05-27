package configuration

const (
	ProfileTypeJira = "JIRA"
)

type JSONConfiguration struct {
	Profiles         []Profile `json:"profiles"`
	LastVersionCheck int64     `json:"lastVersionCheck"`
}

type Profile struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Jira Jira   `json:"jira,omitempty"`
}

type Jira struct {
	Host      string `json:"host"`
	Token     string `json:"token"`
	AccountID string `json:"accountID"`
	Email     string `json:"email"`
	UserKey   string `json:"userKey"`
}
