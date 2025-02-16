package commands

import "fmt"

func HelpCommand(args []string) {
	fmt.Println("Usage: main.go <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  create   Create a new local domain configuration")
	fmt.Println("  list     List all configured local domains")
	fmt.Println("  delete   Delete an existing local domain configuration")
	fmt.Println("  help     Show this help message")
}
