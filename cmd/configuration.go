package cmd

import (
	"bufio"
	Configuration "gira/configuration"
	"os"
	"regexp"
	"strings"

	Services "gira/services"

	"github.com/charmbracelet/x/term"
)

func CmdConfig(configuration *Configuration.Configuration, loggerService *Services.LoggerService) {
	reader := bufio.NewReader(os.Stdin)

	loggerService.Info("Enter the Jira API URL (Example %s):", "https://jira.mycompagny.com")
	if configuration.JSON.JiraHost != "" {
		loggerService.Info("[%s]", configuration.JSON.JiraHost)
	}
	inputJiraHost, _ := reader.ReadString('\n')
	inputJiraHost = strings.TrimSpace(inputJiraHost)
	if inputJiraHost == "" {
		inputJiraHost = configuration.JSON.JiraHost
	}
	if !isValidURLRegex(inputJiraHost) {
		loggerService.Fatal("%s '%s' is not a valid URL. Please make sure it's a full URL including the scheme (e.g. https://example.com)", "ERROR", inputJiraHost)
	}
	configuration.JSON.JiraHost = inputJiraHost

	loggerService.Info("Enter the Jira Token (See %s%s):", inputJiraHost, "/manage-profile/security/api-tokens")
	if configuration.JSON.JiraToken != "" {
		loggerService.Info("[Token already defined. Press %s to continue without making any changes.]", "ENTER")
	}
	inputJiraTokenBytes, _ := term.ReadPassword(os.Stdin.Fd())
	inputJiraToken := strings.TrimSpace(string(inputJiraTokenBytes))
	if inputJiraToken == "" {
		inputJiraToken = configuration.JSON.JiraToken
	}
	configuration.JSON.JiraToken = inputJiraToken

	jiraService := Services.NewJiraService(configuration)
	jiraUser, jiraUserError := jiraService.GetMyself()

	if jiraUserError != nil {
		loggerService.Fatal("❌ Unable to fetch Jira account in %s due to : %v", inputJiraHost, jiraUserError)
	}

	configuration.JSON.JiraEmail = jiraUser.EmailAddress
	configuration.JSON.JiraAccountID = jiraUser.AccountID
	configuration.JSON.JiraUserKey = jiraUser.Key

	configuration.Update()
	loggerService.Info("✅ Done!")
}

func isValidURLRegex(url string) bool {
	re := regexp.MustCompile(`^(https?://)?([\w-]+\.)+[\w-]{2,}(/.*)?$`)
	return re.MatchString(url)
}
