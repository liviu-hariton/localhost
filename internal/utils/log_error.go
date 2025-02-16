package utils

import (
	"fmt"
	"log"
)

func LogError(action string, err error) error {
	formattedErr := fmt.Errorf("%s failed: %s", action, err.Error())
	log.Println(formattedErr) // Logs to stderr for debugging
	return formattedErr       // Return the error to the calling function
}
