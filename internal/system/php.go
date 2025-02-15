package system

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CheckPHPInstalled verifies if PHP is installed on the system.
func CheckPHPInstalled() error {
	cmd := exec.Command("php", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return errors.New("PHP is not installed or not accessible. Install it using Homebrew: 'brew install php'")
	}

	// PHP is installed
	fmt.Println("✔ PHP is installed.")
	return nil
}

// CheckPHPWorking verifies if PHP can execute a basic script.
func CheckPHPWorking() error {
	script := `echo "PHP is working!";`
	cmd := exec.Command("php", "-r", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("PHP is installed, but there was an error running a test script: %s", out.String())
	}

	output := out.String()
	if strings.Contains(output, "PHP is working!") {
		fmt.Println("✔ PHP is working correctly.")
		return nil
	}

	return errors.New("PHP is installed, but it failed the basic test script")
}

// InstallPHP attempts to install PHP using Homebrew.
func InstallPHP() error {
	fmt.Println("PHP is not installed. Attempting to install it using Homebrew...")

	cmd := exec.Command("brew", "install", "php")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install PHP: %s", out.String())
	}

	fmt.Println("✔ PHP installed successfully.")
	return nil
}

// VerifyPHP ensures PHP is installed and working correctly.
func VerifyPHP() error {
	fmt.Println("Checking PHP setup...")

	// Check if PHP is installed
	if err := CheckPHPInstalled(); err != nil {
		fmt.Println(err)

		// Attempt to install PHP
		if installErr := InstallPHP(); installErr != nil {
			return installErr
		}
	}

	// Check if PHP is working
	if err := CheckPHPWorking(); err != nil {
		return err
	}

	return nil
}
