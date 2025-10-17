package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaClient represents an Ollama API client
type OllamaClient struct {
	Host  string
	Model string
}

// GenerateRequest represents the request payload for Ollama API
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse represents the response from Ollama API
type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(host, model string) *OllamaClient {
	return &OllamaClient{
		Host:  host,
		Model: model,
	}
}

// GenerateCommitMessages generates commit message suggestions using Ollama
func (c *OllamaClient) GenerateCommitMessages(diff string, numSuggestions int) ([]string, error) {
	prompt := c.buildPrompt(diff, numSuggestions)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(c.Host+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	var generateResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse the response to extract individual commit messages
	messages := c.parseCommitMessages(generateResp.Response, numSuggestions)

	return messages, nil
}

// buildPrompt creates the prompt for generating commit messages
func (c *OllamaClient) buildPrompt(diff string, numSuggestions int) string {
	return fmt.Sprintf(`You are an expert Git commit message generator. Analyze the following Git diff and generate %d concise, clear commit messages following conventional commit format (feat:, fix:, refactor:, docs:, style:, test:, chore:).

Rules:
- Use conventional commit format
- Be concise but descriptive
- Focus on what changed, not how
- Use present tense ("add feature" not "added feature")
- Keep under 50 characters for the subject line
- Each message should be on a separate line
- Number each message (1., 2., 3., etc.)

Git diff:
%s

Generate %d commit messages:`, numSuggestions, diff, numSuggestions)
}

// parseCommitMessages parses the AI response to extract individual commit messages
func (c *OllamaClient) parseCommitMessages(response string, expectedCount int) []string {
	lines := strings.Split(response, "\n")
	var messages []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove numbering (1., 2., etc.) and clean up
		if len(line) > 2 && (line[1] == '.' || line[2] == '.') {
			// Find the first space after the number
			spaceIndex := strings.Index(line, " ")
			if spaceIndex > 0 {
				line = line[spaceIndex+1:]
			}
		}

		// Clean up any remaining artifacts
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimSpace(line)

		if line != "" {
			messages = append(messages, line)
		}
	}

	// If we got fewer messages than expected, pad with generic ones
	for len(messages) < expectedCount {
		messages = append(messages, fmt.Sprintf("chore: update code (suggestion %d)", len(messages)+1))
	}

	// Return only the requested number
	if len(messages) > expectedCount {
		messages = messages[:expectedCount]
	}

	return messages
}
