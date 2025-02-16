package main

import (
	"fmt"
	"os"

	"github.com/liviu-hariton/localhost/internal/commands"
	"github.com/liviu-hariton/localhost/internal/utils"
)

const Version = "1.0.0"

func main() {
	if utils.HasFlag("--version") {
		fmt.Printf("LocalHost version %s\n", Version)
		return
	}

	// Bypass sudo if dry-run mode is enabled
	if utils.HasFlag("--dry-run") {
		utils.LogInfo("Dry Run mode detected. Skipping privilege escalation.")
	} else {
		// Relaunch the program with sudo if necessary
		if err := utils.RelaunchWithSudo(); err != nil {
			utils.LogError("Relaunch the program with sudo", err)
			os.Exit(1)
		}
	}

	if len(os.Args) < 2 {
		utils.LogWarning("No command provided. Use 'help' for usage information.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		commands.CreateCommand(os.Args[2:])
	case "list":
		commands.ListCommand(os.Args[2:])
	case "delete":
		commands.DeleteCommand(os.Args[2:])
	case "help":
		commands.HelpCommand(os.Args[2:])
	default:
		utils.LogWarning(fmt.Sprintf("Unknown command '%s'. Use 'help' for usage information.\n", os.Args[1]))
		os.Exit(1)
	}
}
