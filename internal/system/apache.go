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

// RestartApache attempts to restart Apache and flushes the DNS cache.
func RestartApache() error {
	if utils.IsDryRun() {
		fmt.Println("DRY RUN: Would restart Apache server and flush the DNS cache.")
		return nil
	}

	// Restart Apache
	restartErr := utils.Spinner("Restarting Apache server...", func() error {
		cmd := exec.Command("sudo", "apachectl", "-k", "restart")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		return cmd.Run()
	})

	if restartErr != nil {
		return fmt.Errorf("failed to restart Apache: %s", restartErr.Error())
	}

	fmt.Println("✔ Apache restarted successfully.")

	if !utils.HasFlag("--no-dns-reset") {
		// Flush DNS cache
		flushErr := utils.Spinner("Flushing DNS cache...", func() error {
			cmd := exec.Command("sudo", "dscacheutil", "-flushcache")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			return cmd.Run()
		})
		if flushErr != nil {
			return fmt.Errorf("failed to flush DNS cache: %s", flushErr.Error())
		}

		// Reset mDNSResponder
		resetErr := utils.Spinner("Resetting mDNSResponder...", func() error {
			cmd := exec.Command("sudo", "killall", "-HUP", "mDNSResponder")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			return cmd.Run()
		})
		if resetErr != nil {
			return fmt.Errorf("failed to reset mDNSResponder: %s", resetErr.Error())
		}

		utils.LogSuccess("DNS cache flushed and mDNSResponder reset successfully.")
	} else {
		utils.LogInfo("Skipping DNS cache flush and mDNSResponder reset as per user request.")
	}

	return nil
}

// VerifyApache checks if Apache is installed and running, and restarts it if needed.
func VerifyApache() error {
	utils.LogInfo("Checking Apache setup...")

	// Check if Apache is installed
	if err := CheckApacheInstalled(); err != nil {
		return err
	}

	// Check if Apache is running
	if err := CheckApacheRunning(); err != nil {
		utils.LogWarning("Apache is not running. Attempting to restart...")

		if restartErr := RestartApache(); restartErr != nil {
			return fmt.Errorf("failed to restart Apache: %s", restartErr.Error())
		}
	}

	return nil
}
