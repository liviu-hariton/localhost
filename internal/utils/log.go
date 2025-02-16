package utils

import (
	"fmt"
	"log"
)

func LogError(action string, err error) error {
	formattedErr := fmt.Errorf("%s[ERROR] %s failed: %s%s", ColorRed, action, err.Error(), ColorReset)
	log.Println(formattedErr) // Logs to stderr for debugging
	return formattedErr       // Return the error to the calling function
}

func LogSuccess(message string) {
	fmt.Printf("%s[SUCCESS] %s%s\n", ColorGreen, message, ColorReset)
}

func LogInfo(message string) {
	fmt.Printf("%s[INFO] %s%s\n", ColorBlue, message, ColorReset)
}

func LogWarning(message string) {
	fmt.Printf("%s[WARNING] %s%s\n", ColorYellow, message, ColorReset)
}

func LogDebug(message string) {
	fmt.Printf("%s[DEBUG] %s%s\n", ColorMagenta, message, ColorReset)
}
