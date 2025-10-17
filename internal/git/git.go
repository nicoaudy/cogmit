package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepo checks if the current directory is a Git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetStagedDiff returns the staged changes diff
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	return string(output), nil
}

// GetWorkingDiff returns the working directory changes diff
func GetWorkingDiff() (string, error) {
	cmd := exec.Command("git", "diff")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get working diff: %w", err)
	}
	return string(output), nil
}

// GetDiff returns staged changes if available, otherwise working directory changes
func GetDiff() (string, error) {
	// First try to get staged changes
	stagedDiff, err := GetStagedDiff()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	// If there are staged changes, return them
	if strings.TrimSpace(stagedDiff) != "" {
		return stagedDiff, nil
	}

	// Otherwise, get working directory changes
	workingDiff, err := GetWorkingDiff()
	if err != nil {
		return "", fmt.Errorf("failed to get working diff: %w", err)
	}

	if strings.TrimSpace(workingDiff) == "" {
		return "", fmt.Errorf("no changes found to commit")
	}

	return workingDiff, nil
}

// HasStagedChanges checks if there are staged changes
func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil // If command fails, there are staged changes
}

// Commit commits with the given message
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}
