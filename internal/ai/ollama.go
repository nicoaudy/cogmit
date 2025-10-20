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
		Timeout: 120 * time.Second, // Increased timeout for larger models
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

func (c *OllamaClient) buildPrompt(diff string, numSuggestions int) string {
	return fmt.Sprintf(`You are an expert software engineer writing clear and professional Git commit messages.

TASK:
Generate %d meaningful commit messages describing the following code diff.

RULES:
- Use the Conventional Commits style (feat:, fix:, refactor:, chore:, etc.)
- Optionally include a scope in parentheses if relevant, e.g. feat(auth):
- Be specific about what changed and why
- Each message must be a single line under 100 characters
- Do NOT repeat or copy any examples
- Return ONLY the commit messages, one per line, no numbering, no explanations

Example style — DO NOT COPY, only follow this style:
feat(auth): add JWT authentication middleware for secure login
fix(api): handle nil pointer exception in user profile update route
refactor(db): move queries to dedicated repository layer for maintainability

DIFF:
%s

Write the commit messages:`, numSuggestions, diff)
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

		// Skip lines that start with explanatory text
		lineLower := strings.ToLower(line)
		if strings.HasPrefix(lineLower, "here are") ||
		   strings.HasPrefix(lineLower, "the following") ||
		   strings.HasPrefix(lineLower, "based on") ||
		   strings.HasPrefix(lineLower, "commit messages") {
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
