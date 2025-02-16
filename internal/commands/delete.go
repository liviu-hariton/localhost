package commands

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/liviu-hariton/localhost/internal/system"
	"github.com/liviu-hariton/localhost/internal/utils"
)

func DeleteCommand(args []string) {
	flagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	domain := flagSet.String("domain", "", "The local domain to delete (e.g., myproject.local)")
	dryRun := flagSet.Bool("dry-run", false, "Simulate changes without modifying any files or directories")
	flagSet.Parse(args)

	// Validate required flags
	if *domain == "" {
		utils.LogWarning("Please provide the -domain flag. For example:")
		fmt.Println("    go run main.go delete -domain=myproject.local")
		flagSet.Usage()
		return
	}

	utils.SetDryRun(*dryRun)
	if utils.IsDryRun() {
		utils.LogInfo(fmt.Sprintf("DRY RUN: Would delete domain '%s' and its references in /etc/hosts.\n", *domain))
		return
	}

	// Ask for user confirmation before proceeding
	if !confirmAction(*domain) {
		utils.LogInfo("Deletion aborted by user.")
		return
	}

	// Remove the virtual host configuration file
	vhostsDir := "/usr/local/etc/httpd/extra/vhosts/"
	vhostFile := vhostsDir + *domain + ".conf"

	if err := os.Remove(vhostFile); err != nil {
		utils.LogError(fmt.Sprintf("Error deleting domain configuration file: %s", err), err)
	} else {
		utils.LogSuccess(fmt.Sprintf("Successfully deleted virtual host configuration for domain '%s'.", *domain))
	}

	// Remove the domain from /etc/hosts
	if err := removeDomainFromHosts(*domain); err != nil {
		utils.LogError(fmt.Sprintf("Error modifying /etc/hosts: %s", err), err)
		return
	}

	utils.LogInfo("Restarting Apache to apply changes...")

	// Restart Apache to apply changes
	if err := system.RestartApache(); err != nil {
		utils.LogError(fmt.Sprintf("Apache Restart Error: %s", err), err)
		return
	}

	utils.LogSuccess(fmt.Sprintf("Successfully removed domain '%s' from /etc/hosts.", *domain))
}

// confirmAction prompts the user for confirmation before proceeding
func confirmAction(domain string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Are you sure you want to delete the domain '%s' and its references in /etc/hosts? (y/N): ", domain)
	response, _ := reader.ReadString('\n')

	// Normalize and trim input
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y"
}

// removeDomainFromHosts removes references to a domain from /etc/hosts
func removeDomainFromHosts(domain string) error {
	hostsFilePath := "/etc/hosts"

	// Open the /etc/hosts file for reading
	file, err := os.Open(hostsFilePath)
	if err != nil {
		return fmt.Errorf("failed to open /etc/hosts: %s", err.Error())
	}
	defer file.Close()

	var newLines []string
	scanner := bufio.NewScanner(file)

	// Iterate through each line and exclude lines referencing the domain
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore lines that reference the domain
		if strings.Contains(line, domain) {
			utils.LogInfo(fmt.Sprintf("Removing line from /etc/hosts: %s\n", line))
			continue
		}

		newLines = append(newLines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read /etc/hosts: %s", err.Error())
	}

	// Write the updated lines back to the file
	if err := writeLinesToFile(newLines, hostsFilePath); err != nil {
		return fmt.Errorf("failed to update /etc/hosts: %s", err.Error())
	}

	return nil
}

// writeLinesToFile writes a slice of strings to a file, overwriting its contents
func writeLinesToFile(lines []string, filePath string) error {
	// Use a buffer to prepare the new content
	var buffer bytes.Buffer
	for _, line := range lines {
		buffer.WriteString(line + "\n")
	}

	// Use sudo to write to /etc/hosts without printing the output
	cmd := exec.Command("sudo", "tee", filePath)
	cmd.Stdin = &buffer
	cmd.Stdout = nil // Suppress output to the screen
	cmd.Stderr = nil // Suppress error output to the screen

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write to %s: %s", filePath, err.Error())
	}

	return nil
}
