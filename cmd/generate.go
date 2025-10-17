package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nicoaudy/cogmit/internal/ai"
	"github.com/nicoaudy/cogmit/internal/config"
	"github.com/nicoaudy/cogmit/internal/git"
	"github.com/nicoaudy/cogmit/internal/ui"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate commit messages for current changes",
	Long: `Analyze your Git changes and generate intelligent commit message suggestions
using your configured Ollama model.

The tool will:
1. Check if you're in a Git repository
2. Analyze your staged or working directory changes
3. Generate commit message suggestions using AI
4. Let you select or edit a message
5. Optionally commit automatically based on your configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runGenerate(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runGenerate() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if we're in a Git repository
	if !git.IsGitRepo() {
		return fmt.Errorf("not in a Git repository")
	}

	// Get the diff
	diff, err := git.GetDiff()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("no changes found to commit")
	}

	// Generate commit messages using Ollama
	ollamaClient := ai.NewOllamaClient(cfg.OllamaHost, cfg.Model)
	messages, err := ollamaClient.GenerateCommitMessages(diff, cfg.NumSuggestions)
	if err != nil {
		// Provide more helpful error messages
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			return fmt.Errorf("request timed out - the model may be too slow. Try a smaller model like 'llama3.2:1b' or 'llama3.2:3b'")
		}
		if strings.Contains(err.Error(), "connection refused") {
			return fmt.Errorf("cannot connect to Ollama - make sure it's running with 'ollama serve'")
		}
		return fmt.Errorf("failed to generate commit messages: %w", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("no commit messages generated")
	}

	// Show interactive selector
	selector := ui.NewSelectorModel(messages)
	program := tea.NewProgram(selector)

	finalModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to run selector: %w", err)
	}

	selectorModel := finalModel.(ui.SelectorModel)
	selectedMessage := selectorModel.GetSelected()

	if selectedMessage == "" {
		fmt.Println("❌ No commit message selected")
		return nil
	}

	// Ask for confirmation if not auto-committing
	if !cfg.AutoCommit {
		fmt.Print("Commit with this message? [Y/n]: ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "n" || response == "no" {
			fmt.Println("❌ Commit cancelled")
			return nil
		}
	}

	// Commit the changes
	if err := git.Commit(selectedMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
