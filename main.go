package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/liviu-hariton/localhost/internal/config"
	"github.com/liviu-hariton/localhost/internal/system"
	"github.com/liviu-hariton/localhost/internal/utils"
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Simulate changes without modifying any files or directories")
}

func main() {
	// Relaunch with sudo if necessary
	if err := utils.RelaunchWithSudo(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// Define command-line flags
	domain := flag.String("domain", "", "The local domain to set up (e.g., myproject.local)")
	doc_root := flag.String("doc_root", "", "The document root for the virtual host")

	flag.Parse()

	// Pass the dryRun mode to the system package
	utils.SetDryRun(dryRun)

	if dryRun {
		fmt.Println("Running in Dry Run mode: No changes will be made.")
	}

	// Ensure the domain and document root are provided
	if *domain == "" || *doc_root == "" {
		fmt.Println("Error: Please provide both -domain and -docRoot flags. For example:")
		fmt.Println("    go run main.go -domain=myproject.local -doc_root=/path/on/disk/to/myproject")
		os.Exit(1)
	}

	fmt.Printf("Starting setup for domain: %s\n", *domain)

	fmt.Println("Starting system checks...")

	// Check Apache
	if err := system.VerifyApache(); err != nil {
		fmt.Printf("Apache Error: %s\n", err)
		return
	}

	// Check MySQL
	if err := system.VerifyMySQL(); err != nil {
		fmt.Printf("MySQL Error: %s\n", err)
		return
	}

	// Check PHP
	if err := system.VerifyPHP(); err != nil {
		fmt.Printf("PHP Error: %s\n", err)
		return
	}

	fmt.Println("All checks passed successfully!")

	// Modify Hosts File
	if err := config.ModifyHosts(*domain); err != nil {
		fmt.Printf("Hosts File Error: %s\n", err)
		return
	}

	fmt.Println("Ensuring vhosts are enabled and adding virtual host...")

	// Ensure vhosts are enabled
	if err := config.EnsureVhostsEnabled(); err != nil {
		fmt.Printf("Apache Config Error: %s\n", err)
		return
	}

	// Add Virtual Host
	if err := config.AddVirtualHost(*domain, *doc_root); err != nil {
		fmt.Printf("Virtual Host Error: %s\n", err)
		return
	}

	fmt.Println("Restarting Apache to apply changes...")

	// Restart Apache to apply changes
	if err := system.RestartApache(); err != nil {
		fmt.Printf("Apache Restart Error: %s\n", err)
		return
	}

	// Ensure SSL Certificates
	if err := system.EnsureSSLCertificates(); err != nil {
		fmt.Printf("SSL Error: %s\n", err)
		return
	}

	fmt.Println("All changes applied successfully!")

	fmt.Printf("You should be able now to access your new project at http://%s or https://%s\n", *domain, *domain)
}
