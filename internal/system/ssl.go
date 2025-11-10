package system

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/liviu-hariton/localhost/internal/utils"
)

// Default SSL certificate paths
const sslCertificateFile = "/usr/local/etc/httpd/ssl/server.crt"
const sslCertificateKeyFile = "/usr/local/etc/httpd/ssl/server.key"

// EnsureSSLCertificates ensures the SSL certificate and key files exist, generating them if necessary.
func EnsureSSLCertificates() error {
	utils.LogInfo("Checking for SSL certificates...")

	// Check if the certificate and key files exist
	if _, err := os.Stat(sslCertificateFile); os.IsNotExist(err) {
		utils.LogWarning(fmt.Sprintf("SSL certificate not found at %s. Generating a self-signed certificate...\n", sslCertificateFile))

		// Ensure the /usr/local/etc/httpd/ssl directory exists
		if !utils.IsDryRun() {
			if err := utils.CreateDirectory("/usr/local/etc/httpd/ssl"); err != nil {
				return utils.LogError("Creating SSL directory", err)
			}
		} else {
			utils.LogInfo("DRY RUN: Would create the SSL directory.")
		}

		// Generate a self-signed certificate
		if !utils.IsDryRun() {
			cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-days", "365", "-newkey", "rsa:2048",
				"-keyout", sslCertificateKeyFile,
				"-out", sslCertificateFile,
				"-subj", "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return utils.LogError("Generating self-signed SSL certificate", err)
			}

			utils.LogSuccess(fmt.Sprintf("✔ Self-signed SSL certificate generated:\n - Certificate: %s\n - Key: %s\n", sslCertificateFile, sslCertificateKeyFile))
		} else {
			utils.LogInfo("DRY RUN: Would Generate a self-signed certificate.")
		}

	} else {
		utils.LogSuccess("SSL certificates already exist.")
	}

	return nil
}
