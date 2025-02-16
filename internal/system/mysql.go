package system

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/liviu-hariton/localhost/internal/utils"
)

// CheckMySQLInstalled verifies if MySQL is installed on the system.
func CheckMySQLInstalled() error {
	cmd := exec.Command("mysql", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return errors.New("MySQL is not installed or not accessible. Install it using Homebrew: 'brew install mysql'")
	}

	// MySQL is installed
	fmt.Println("✔ MySQL is installed.")
	return nil
}

// CheckMySQLRunning verifies if MySQL is currently running.
func CheckMySQLRunning() error {
	cmd := exec.Command("brew", "services", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to check running services: %s", err.Error())
	}

	if strings.Contains(out.String(), "mysql") && strings.Contains(out.String(), "started") {
		fmt.Println("✔ MySQL is running.")
		return nil
	}

	return errors.New("MySQL is not running")
}

// InstallMySQL attempts to install MySQL using Homebrew.
func InstallMySQL() error {
	if utils.IsDryRun() {
		utils.LogInfo("DRY RUN: Would install MySQL using Homebrew.")
		return nil
	}

	utils.LogWarning("MySQL is not installed. Attempting to install it using Homebrew...")

	cmd := exec.Command("brew", "install", "mysql")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install MySQL: %s", out.String())
	}

	utils.LogSuccess("MySQL installed successfully.")
	return nil
}

// RestartMySQL attempts to restart MySQL using Homebrew services.
func RestartMySQL() error {
	if utils.IsDryRun() {
		utils.LogInfo("DRY RUN: Would restart MySQL.")
		return nil
	}

	cmd := exec.Command("brew", "services", "restart", "mysql")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to restart MySQL: %s", out.String())
	}

	utils.LogSuccess("MySQL restarted successfully.")
	return nil
}

// VerifyMySQL ensures MySQL is installed, running, and starts it if needed.
func VerifyMySQL() error {
	utils.LogInfo("Checking MySQL setup...")

	// Check if MySQL is installed
	if err := CheckMySQLInstalled(); err != nil {
		fmt.Println(err)

		// Attempt to install MySQL
		if installErr := InstallMySQL(); installErr != nil {
			return installErr
		}
	}

	// Check if MySQL is running
	if err := CheckMySQLRunning(); err != nil {
		utils.LogWarning("MySQL is not running. Attempting to restart...")

		if restartErr := RestartMySQL(); restartErr != nil {
			return fmt.Errorf("failed to restart MySQL: %s", restartErr.Error())
		}
	}

	return nil
}
