package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cogmit",
	Short: "AI-powered Git commit message generator using local Ollama models",
	Long: `cogmit is a CLI tool that generates smart Git commit messages using local AI models via Ollama.

It analyzes your staged changes and provides intelligent commit message suggestions
that you can select from using an interactive interface.

Examples:
  cogmit          # Generate commit messages for staged changes
  cogmit setup    # Configure Ollama settings and preferences`,
	Run: func(cmd *cobra.Command, args []string) {
		// This will be handled by the generate command
		generateCmd.Run(cmd, args)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(generateCmd)
}
