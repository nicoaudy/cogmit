package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nicoaudy/cogmit/internal/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure cogmit settings",
	Long: `Setup cogmit by configuring your Ollama host, model preferences, and other settings.

This will create a configuration file at ~/.config/cogmit/config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSetup(); err != nil {
			fmt.Printf("❌ Setup failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func runSetup() error {
	fmt.Println("🔧 Setting up cogmit...")
	fmt.Println()

	// Load existing config or use defaults
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not load existing config: %v\n", err)
		cfg = config.DefaultConfig()
	}

	reader := bufio.NewReader(os.Stdin)

	// Ollama host
	fmt.Printf("Ollama host [%s]: ", cfg.OllamaHost)
	ollamaHost, _ := reader.ReadString('\n')
	ollamaHost = strings.TrimSpace(ollamaHost)
	if ollamaHost != "" {
		cfg.OllamaHost = ollamaHost
	}

	// Model name
	fmt.Printf("Model name [%s]: ", cfg.Model)
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model != "" {
		cfg.Model = model
	}

	// Number of suggestions
	fmt.Printf("Number of suggestions [%d]: ", cfg.NumSuggestions)
	numSuggestionsStr, _ := reader.ReadString('\n')
	numSuggestionsStr = strings.TrimSpace(numSuggestionsStr)
	if numSuggestionsStr != "" {
		if num, err := strconv.Atoi(numSuggestionsStr); err == nil && num > 0 {
			cfg.NumSuggestions = num
		}
	}

	// Auto commit
	fmt.Printf("Auto-commit after selection? [y/N]: ")
	autoCommitStr, _ := reader.ReadString('\n')
	autoCommitStr = strings.TrimSpace(strings.ToLower(autoCommitStr))
	cfg.AutoCommit = autoCommitStr == "y" || autoCommitStr == "yes"

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("📁 Config location: ~/.config/cogmit/config.json\n")
	fmt.Println()
	fmt.Println("🎉 Setup complete! You can now run 'cogmit' to generate commit messages.")

	return nil
}
