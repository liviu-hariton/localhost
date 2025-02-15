package system

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HttpdConfPath defines the path to the Apache main configuration file.
const HttpdConfPath = "/usr/local/etc/httpd/httpd.conf"

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

		// After installing PHP, update Apache configuration
		if err := EnablePHPModuleInHttpdConf(); err != nil {
			return err
		}

		// Restart Apache to apply changes
		if restartErr := RestartApache(); restartErr != nil {
			return fmt.Errorf("failed to restart Apache after enabling PHP: %s", restartErr.Error())
		}
	}

	// Check if PHP is working
	if err := CheckPHPWorking(); err != nil {
		return err
	}

	return nil
}

func EnablePHPModuleInHttpdConf() error {
	fmt.Println("Enabling PHP module in Apache configuration...")

	file, err := os.Open(HttpdConfPath)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf: %s", err.Error())
	}
	defer file.Close()

	var lines []string
	phpModuleLoaded := false
	setHandlerAdded := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the PHP module is already loaded
		if strings.Contains(line, "LoadModule php") && !strings.HasPrefix(line, "#") {
			phpModuleLoaded = true
		}

		// Check if the SetHandler directive already exists
		if strings.Contains(line, "SetHandler application/x-httpd-php") {
			setHandlerAdded = true
		}

		// Uncomment the PHP module loading line if it's commented
		if strings.Contains(line, "LoadModule php") && strings.HasPrefix(line, "#") {
			line = strings.TrimPrefix(line, "#")
			phpModuleLoaded = true
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read httpd.conf: %s", err.Error())
	}

	// Add the PHP module loading line if not found
	if !phpModuleLoaded {
		lines = append([]string{"LoadModule php_module /usr/local/opt/php/lib/httpd/modules/libphp.so"}, lines...)
		fmt.Println("✔ Added PHP module loading line to httpd.conf.")
	}

	// Add the SetHandler directive if not found
	if !setHandlerAdded {
		lines = append(lines, "<FilesMatch \\.php$>")
		lines = append(lines, "    SetHandler application/x-httpd-php")
		lines = append(lines, "</FilesMatch>")
		fmt.Println("✔ Added SetHandler directive to httpd.conf.")
	}

	// Write the updated content back to httpd.conf
	file, err = os.OpenFile(HttpdConfPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf for writing: %s", err.Error())
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to httpd.conf: %s", err.Error())
		}
	}

	fmt.Println("✔ PHP module and handler enabled in httpd.conf.")
	return nil
}
