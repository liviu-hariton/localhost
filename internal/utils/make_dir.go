package utils

import (
	"fmt"
	"os"
)

func CreateDirectory(path string) error {
	if IsDryRun() {
		fmt.Printf("DRY RUN: Would create directory: %s\n", path)
		return nil
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return LogError(fmt.Sprintf("Creating directory %s", path), err)
	}

	fmt.Printf("✔ Created directory: %s\n", path)
	return nil
}
