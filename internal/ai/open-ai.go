package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"

	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type OpenAI struct {
	logger *log.Logger
	client openai.Client
	Model  string
}

func NewOpenAI(logger *log.Logger) *OpenAI {
	endpoint := os.Getenv("GIRA_AI_ENDPOINT")
	apikey := os.Getenv("GIRA_AI_APIKEY")
	model := os.Getenv("GIRA_AI_MODEL")
	if endpoint == "" || model == "" {
		logger.Fatal("❌ AI configuration error: missing required environment variables. See", "https://github.com/Ealenn/gira")
	}

	configuration := []option.RequestOption{}
	configuration = append(configuration, option.WithBaseURL(endpoint))
	if apikey != "" {
		configuration = append(configuration, option.WithAPIKey(apikey))
	}

	return &OpenAI{
		logger: logger,
		Model:  model,
		client: openai.NewClient(configuration...),
	}
}

func (agent *OpenAI) BranchNames(issue *issue.Issue) ([]string, error) {
	prompt := fmt.Sprintf(
		"Based on this Ticket:\n"+
			"Title: %s\nDescription: %s\n\n"+
			"Generate exactly 3 concise git branch names that are:\n"+
			"- lowercase\n"+
			"- hyphen-separated\n"+
			"- composed only of words from the title/description, without adjectives like 'quick', 'new', 'urgent', 'bug', 'feature'\n"+
			"- do not add prefixes like 'feat/', 'fix/', or 'branch/'\n"+
			"- keep names strictly relevant and descriptive of the task\n"+
			"Return **only a JSON array** like this:\n"+
			"[\"example-branch-1\", \"example-branch-2\", \"example-branch-3\"]",
		issue.Title, agent.getShortIssueDescription(issue),
	)

	respose, err := agent.askJSONStringArray(prompt)
	if err != nil {
		agent.logger.Debug("Unable to generate branch name on model %s due to %v", agent.Model, err)
		return nil, err
	}

	return respose, nil
}

func (agent *OpenAI) CommitNames(issue *issue.Issue) ([]string, error) {
	prompt := fmt.Sprintf(
		"Based on this Ticket:\n"+
			"Title: %s\nDescription: %s\n\n"+
			"Generate exactly 3 concise git commit messages following the Conventional Commits specification, based on this ticket.\n"+
			"- Allowed types: feat, fix, docs, style, refactor, perf, test, chore.\n"+
			"- Scope is optional but must be lowercase if present.\n"+
			"- Message should be short, imperative, and descriptive.\n"+
			"Return ONLY a valid JSON array of strings, not markdown, e.g. [\"fix(auth): resolve login bug after password reset\", \"feat: improve session handling\", \"chore: update dependencies\"].",
		issue.Title, agent.getShortIssueDescription(issue),
	)
	return agent.askJSONStringArray(prompt)
}

func (agent *OpenAI) IssueSummary(issue *issue.Issue) (string, error) {
	prompt := fmt.Sprintf(
		"Based on this Ticket:\n"+
			"Assignees: %v\nStatus: %s\nTypes: %v\nTitle: %s\nDescription: %s\n\n"+
			"Generate concise summary.\n"+
			"Return ONLY text format",
		issue.Assignees, issue.Status, issue.Types, issue.Title, agent.getShortIssueDescription(issue),
	)
	return agent.askString(prompt)
}

func (agent *OpenAI) Rewrite(context string, text string) (string, error) {
	prompt := fmt.Sprintf(
		"Context: %s\n\nRewrite the text below to be clearer and more precise, without adding extra information or changing meaning:\n\n%s",
		context, text,
	)
	return agent.askString(prompt)
}

func (agent *OpenAI) getShortIssueDescription(issue *issue.Issue) string {
	description := issue.Description
	if len(description) > 4096 {
		description = description[:4096]
	}

	return description
}

func (agent *OpenAI) askString(prompt string) (string, error) {
	systemMessage := openai.SystemMessage("You are a git workflow assistant. Respond **only** with a text. Do not include markdown, backticks, or any other text.")
	chatCompletion, err := agent.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			systemMessage,
			openai.UserMessage(prompt),
		},
		Model: agent.Model,
	})

	if err != nil {
		agent.logger.Debug("OpenAI error %s", err)
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		agent.logger.Debug("Error: response body contains no choices")
		return "", fmt.Errorf("no response from AI API")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (agent *OpenAI) askJSONStringArray(prompt string) ([]string, error) {
	systemMessage := openai.SystemMessage("You are a git workflow assistant. Respond **only** with a valid JSON array of strings. Do not include markdown, backticks, or any other text.")
	chatCompletion, err := agent.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			systemMessage,
			openai.UserMessage(prompt),
		},
		Model: agent.Model,
	})

	if err != nil {
		agent.logger.Debug("OpenAI error %s", err)
		return nil, err
	}

	if len(chatCompletion.Choices) == 0 {
		agent.logger.Debug("Error: response body contains no choices")
		return nil, fmt.Errorf("no response from AI API")
	}

	var output []string
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &output)
	if err != nil {
		agent.logger.Debug("Error: model output response is not valid JSON")
		return nil, fmt.Errorf("failed to parse JSON array from model output: %v", err)
	}

	return output, nil
}
