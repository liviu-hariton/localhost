package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/liviu-hariton/localhost/internal/config"
	"github.com/liviu-hariton/localhost/internal/system"
	"github.com/liviu-hariton/localhost/internal/utils"
)

func main() {
	// Relaunch with sudo if necessary
	if err := utils.RelaunchWithSudo(); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// Define a command-line flag for the domain name
	domain := flag.String("domain", "", "The local domain to set up (e.g., myproject.local)")
	flag.Parse()

	// Ensure the domain is provided
	if *domain == "" {
		fmt.Println("Error: Please provide a domain name using the -domain flag. For example:")
		fmt.Println("    go run main.go -domain=myproject.local")
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

	// Modify Hosts File
	if err := config.ModifyHosts(*domain); err != nil {
		fmt.Printf("Hosts File Error: %s\n", err)
		return
	}

	fmt.Println("All checks passed successfully!")
}
