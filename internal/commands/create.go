package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/liviu-hariton/localhost/internal/config"
	"github.com/liviu-hariton/localhost/internal/system"
	"github.com/liviu-hariton/localhost/internal/utils"
)

func CreateCommand(args []string) {
	// Define command-line flags for the create subcommand
	flagSet := flag.NewFlagSet("create", flag.ExitOnError)
	domain := flagSet.String("domain", "", "The local domain to set up (e.g., myproject.local)")
	docRoot := flagSet.String("doc_root", "", "The document root for the virtual host")
	dryRun := flagSet.Bool("dry-run", false, "Simulate changes without modifying any files or directories")
	flagSet.Parse(args)

	// Validate required flags
	if *domain == "" || *docRoot == "" {
		utils.LogWarning("Please provide both -domain and -doc_root flags. For example:")
		utils.LogWarning("    go run main.go create -domain=myproject.local -doc_root=/path/on/disk/to/myproject")
		os.Exit(1)
	}

	// Set dry run mode
	utils.SetDryRun(*dryRun)
	if *dryRun {
		utils.LogInfo("Running in Dry Run mode: No changes will be made.")
	}

	utils.LogInfo(fmt.Sprintf("Starting setup for domain: %s\n", *domain))

	fmt.Println("Starting system checks...")

	// Check Apache
	if err := system.VerifyApache(); err != nil {
		utils.LogError(fmt.Sprintf("Apache Error: %s\n", err), err)
		return
	}

	// Check MySQL
	if err := system.VerifyMySQL(); err != nil {
		utils.LogError(fmt.Sprintf("MySQL Error: %s\n", err), err)
		return
	}

	// Check PHP
	if err := system.VerifyPHP(); err != nil {
		utils.LogError(fmt.Sprintf("PHP Error: %s\n", err), err)
		return
	}

	utils.LogSuccess("All checks passed successfully!")

	// Modify Hosts File
	if err := config.ModifyHosts(*domain); err != nil {
		utils.LogError(fmt.Sprintf("Hosts File Error: %s\n", err), err)
		return
	}

	fmt.Println("Ensuring vhosts are enabled and adding virtual host...")

	// Ensure vhosts are enabled
	if err := config.EnsureVhostsEnabled(); err != nil {
		utils.LogError(fmt.Sprintf("Apache Config Error: %s\n", err), err)
		return
	}

	// Add Virtual Host
	if err := config.AddVirtualHost(*domain, *docRoot); err != nil {
		utils.LogError(fmt.Sprintf("Virtual Host Error: %s\n", err), err)
		return
	}

	utils.LogInfo("Restarting Apache to apply changes...")

	// Restart Apache to apply changes
	if err := system.RestartApache(); err != nil {
		utils.LogError(fmt.Sprintf("Apache Restart Error: %s\n", err), err)
		return
	}

	// Ensure SSL Certificates
	if err := system.EnsureSSLCertificates(); err != nil {
		utils.LogError(fmt.Sprintf("SSL Error: %s\n", err), err)
		return
	}

	utils.LogSuccess("All changes applied successfully!")

	utils.LogInfo(fmt.Sprintf("You should now be able to access your new project at http://%s or https://%s\n", *domain, *domain))
}
