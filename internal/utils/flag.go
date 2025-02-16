package utils

import "os"

// HasFlag checks if a specific flag is present in the command-line arguments
func HasFlag(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}
