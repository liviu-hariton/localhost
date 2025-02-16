package system

import (
	"fmt"
	"os"
	"os/exec"
)

// Default SSL certificate paths
const sslCertificateFile = "/etc/apache2/ssl/server.crt"
const sslCertificateKeyFile = "/etc/apache2/ssl/server.key"

// EnsureSSLCertificates ensures the SSL certificate and key files exist, generating them if necessary.
func EnsureSSLCertificates() error {
	fmt.Println("Checking for SSL certificates...")

	// Check if the certificate and key files exist
	if _, err := os.Stat(sslCertificateFile); os.IsNotExist(err) {
		fmt.Printf("SSL certificate not found at %s. Generating a self-signed certificate...\n", sslCertificateFile)

		// Ensure the /etc/apache2/ssl directory exists
		if err := os.MkdirAll("/etc/apache2/ssl", 0755); err != nil {
			return fmt.Errorf("failed to create SSL directory: %s", err.Error())
		}

		// Generate a self-signed certificate
		cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-days", "365", "-newkey", "rsa:2048",
			"-keyout", sslCertificateKeyFile,
			"-out", sslCertificateFile,
			"-subj", "/C=US/ST=State/L=City/O=Organization/OU=Unit/CN=localhost")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to generate self-signed SSL certificate: %s", err.Error())
		}

		fmt.Printf("✔ Self-signed SSL certificate generated:\n - Certificate: %s\n - Key: %s\n", sslCertificateFile, sslCertificateKeyFile)
	} else {
		fmt.Println("✔ SSL certificates already exist.")
	}

	return nil
}
