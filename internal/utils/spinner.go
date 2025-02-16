package utils

import (
	"fmt"
	"time"
)

// Spinner displays a rotating spinner until the provided function completes.
func Spinner(message string, fn func() error) error {
	stop := make(chan bool)
	errCh := make(chan error)

	// Start the spinner in a goroutine
	go func() {
		frames := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				fmt.Printf("\r%s %s", message, frames[i])
				i = (i + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Execute the function in the main thread
	err := fn()

	// Stop the spinner and clear the line
	stop <- true
	fmt.Printf("\r%s ✔\n", message)

	// Pass along any error from the function
	if err != nil {
		errCh <- err
		close(errCh)
		return err
	}

	close(errCh)
	return nil
}
