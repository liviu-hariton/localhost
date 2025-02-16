package utils

var dryRun bool // Global variable to track dry run mode

// SetDryRun sets the dry run mode for the utility
func SetDryRun(mode bool) {
	dryRun = mode
}

// IsDryRun checks if the utility is running in dry run mode
func IsDryRun() bool {
	return dryRun
}
