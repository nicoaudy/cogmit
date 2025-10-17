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
	fmt.Println("🔍 Analyzing changes...")
	diff, err := git.GetDiff()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("no changes found to commit")
	}

	// Check if there are staged changes
	hasStaged := git.HasStagedChanges()
	if hasStaged {
		fmt.Println("📝 Found staged changes")
	} else {
		fmt.Println("📝 Found working directory changes")
	}

	// Generate commit messages using Ollama
	fmt.Printf("🤖 Generating commit messages using %s...\n", cfg.Model)

	ollamaClient := ai.NewOllamaClient(cfg.OllamaHost, cfg.Model)
	messages, err := ollamaClient.GenerateCommitMessages(diff, cfg.NumSuggestions)
	if err != nil {
		return fmt.Errorf("failed to generate commit messages: %w", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("no commit messages generated")
	}

	fmt.Printf("✨ Generated %d commit message suggestions\n\n", len(messages))

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

	// Show the selected message
	fmt.Printf("\n✅ Selected: %s\n", selectedMessage)

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
	fmt.Println("💾 Committing changes...")
	if err := git.Commit(selectedMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	fmt.Println("🎉 Successfully committed!")
	return nil
}
