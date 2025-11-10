package commands

import (
	"fmt"
	"os"

	"github.com/liviu-hariton/localhost/internal/utils"
)

func ListCommand(args []string) {
	vhostsDir := "/opt/homebrew/etc/httpd/extra/vhosts/"
	files, err := os.ReadDir(vhostsDir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Error reading vhosts directory: %s\n", err), err)
		return
	}

	fmt.Println("Configured domains:")
	for _, file := range files {
		if file.IsDir() || !file.Type().IsRegular() || file.Name() == ".DS_Store" {
			continue
		}
		utils.LogDebug(file.Name())
	}
}
