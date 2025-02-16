/**
* Verify that Apache is installed and running. Restart it if not running.
 */

package system

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/liviu-hariton/localhost/internal/utils"
)

// CheckApacheInstalled verifies if Apache is installed on the system.
func CheckApacheInstalled() error {
	cmd := exec.Command("apachectl", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("apache is not installed or not accessible: %s", out.String())
	}

	// Apache is installed
	fmt.Println("✔ Apache is installed.")
	return nil
}

// CheckApacheRunning verifies if Apache is currently running.
func CheckApacheRunning() error {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to check running processes: %s", err.Error())
	}

	if strings.Contains(out.String(), "httpd") {
		fmt.Println("✔ Apache is running.")
		return nil
	}

	return errors.New("apache is not running")
}

// RestartApache attempts to restart Apache if it is not running.
func RestartApache() error {
	if utils.IsDryRun() {
		fmt.Println("DRY RUN: Would restart Apache server.")
		return nil
	}

	cmd := exec.Command("sudo", "apachectl", "-k", "restart")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to restart Apache: %s", out.String())
	}

	fmt.Println("✔ Apache restarted successfully.")
	return nil
}

// VerifyApache checks if Apache is installed and running, and restarts it if needed.
func VerifyApache() error {
	fmt.Println("Checking Apache setup...")

	// Check if Apache is installed
	if err := CheckApacheInstalled(); err != nil {
		return err
	}

	// Check if Apache is running
	if err := CheckApacheRunning(); err != nil {
		fmt.Println("Apache is not running. Attempting to restart...")

		if restartErr := RestartApache(); restartErr != nil {
			return fmt.Errorf("failed to restart Apache: %s", restartErr.Error())
		}
	}

	return nil
}
