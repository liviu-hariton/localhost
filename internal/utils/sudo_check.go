package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// RelaunchWithSudo relaunches the current program with sudo if permissions are insufficient.
func RelaunchWithSudo() error {
	// Check if the program is already running with sudo
	if os.Geteuid() == 0 {
		// Already running as root, no need to relaunch
		return nil
	}

	// Relaunch the program with sudo
	fmt.Println("Insufficient permissions. Relaunching with sudo...")

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
