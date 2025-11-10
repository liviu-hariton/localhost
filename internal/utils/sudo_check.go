package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RelaunchWithSudo relaunches the current program with sudo if permissions are insufficient.
func RelaunchWithSudo() error {
	// Check if the program is already running with sudo
	if os.Geteuid() == 0 {
		// Already running as root, no need to relaunch
		return nil
	}

	// Relaunch the program with sudo
	LogWarning("Insufficient permissions. Relaunching with sudo...")

	cmd := exec.Command("sudo", append([]string{os.Args[0]}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and replace the current process
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run with sudo: %w", err)
	}

	// Exit the current process
	os.Exit(0)
	return nil // This line will never be reached
}

// GetOriginalUser returns the original username when running with sudo, or the current user otherwise.
func GetOriginalUser() string {
	// When running with sudo, SUDO_USER contains the original username
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		return sudoUser
	}
	// Otherwise, return the current user
	return os.Getenv("USER")
}

// RunAsOriginalUser runs a command as the original user (not root).
// This is necessary for Homebrew commands which refuse to run as root.
func RunAsOriginalUser(cmd *exec.Cmd) error {
	// If we're not running as root, just run the command normally
	if os.Geteuid() != 0 {
		return cmd.Run()
	}

	// Get the original user
	originalUser := GetOriginalUser()
	if originalUser == "" {
		return fmt.Errorf("could not determine original user")
	}

	// Get the original user's home directory
	originalHome := os.Getenv("SUDO_HOME")
	if originalHome == "" {
		// Fallback: construct home path (works on macOS/Linux)
		originalHome = fmt.Sprintf("/Users/%s", originalUser)
	}

	// Build the command to run as the original user
	// We need to use sudo -u to switch to the original user
	// Use -E to preserve environment variables, but we'll override HOME
	args := append([]string{"-u", originalUser}, cmd.Args...)
	sudoCmd := exec.Command("sudo", args...)
	sudoCmd.Stdin = cmd.Stdin
	sudoCmd.Stdout = cmd.Stdout
	sudoCmd.Stderr = cmd.Stderr

	// Set up environment: preserve existing env vars but override HOME and USER
	env := os.Environ()
	// Remove existing HOME and USER entries and add new ones
	newEnv := make([]string, 0, len(env)+2)
	for _, e := range env {
		if !strings.HasPrefix(e, "HOME=") && !strings.HasPrefix(e, "USER=") {
			newEnv = append(newEnv, e)
		}
	}
	newEnv = append(newEnv, fmt.Sprintf("HOME=%s", originalHome))
	newEnv = append(newEnv, fmt.Sprintf("USER=%s", originalUser))
	// Preserve any env vars from the original command
	if cmd.Env != nil {
		newEnv = append(newEnv, cmd.Env...)
	}
	sudoCmd.Env = newEnv

	return sudoCmd.Run()
}
