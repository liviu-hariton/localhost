package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HostsFilePath defines the path to the hosts file
const HostsFilePath = "/etc/hosts"

// CheckDomainInHosts checks if the domain already exists in the hosts file.
func CheckDomainInHosts(domain string) (bool, error) {
	file, err := os.Open(HostsFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to open hosts file: %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, domain) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to read hosts file: %s", err.Error())
	}

	return false, nil
}

// AddDomainToHosts adds the domain to the hosts file if it doesn't already exist.
func AddDomainToHosts(domain string) error {
	exists, err := CheckDomainInHosts(domain)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("✔ The domain '%s' already exists in the hosts file.\n", domain)
		return nil
	}

	// Open the hosts file for appending
	file, err := os.OpenFile(HostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: you must run this program with elevated permissions (e.g., using sudo)")
		}
		return fmt.Errorf("failed to open hosts file for writing: %s", err.Error())
	}
	defer file.Close()

	// Add the new domain
	newEntry := fmt.Sprintf("\n127.0.0.1 %s\n", domain)
	if _, err := file.WriteString(newEntry); err != nil {
		return fmt.Errorf("failed to write to hosts file: %s", err.Error())
	}

	fmt.Printf("✔ Successfully added '%s' to the hosts file.\n", domain)
	return nil
}

// ModifyHosts handles the entire process of adding a domain to the hosts file.
func ModifyHosts(domain string) error {
	fmt.Printf("Modifying hosts file for domain: %s\n", domain)

	err := AddDomainToHosts(domain)
	if err != nil {
		return fmt.Errorf("failed to modify hosts file: %s", err.Error())
	}

	return nil
}
