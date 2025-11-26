package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/liviu-hariton/localhost/internal/utils"
)

// Default SSL certificate paths
const sslCertificateFile = "/opt/homebrew/etc/httpd/ssl/server.crt"
const sslCertificateKeyFile = "/opt/homebrew/etc/httpd/ssl/server.key"

// HttpdSSLConfPath defines the path to the Apache SSL configuration file.
const HttpdSSLConfPath = "/opt/homebrew/etc/httpd/extra/httpd-ssl.conf"

// EnsureSSLCertificates ensures the SSL certificate and key files exist, generating them if necessary.
func EnsureSSLCertificates() error {
	utils.LogInfo("Checking for SSL certificates...")

	// Check if the certificate and key files exist
	if _, err := os.Stat(sslCertificateFile); os.IsNotExist(err) {
		utils.LogWarning(fmt.Sprintf("SSL certificate not found at %s. Generating a self-signed certificate...\n", sslCertificateFile))

		// Ensure the /opt/homebrew/etc/httpd/ssl directory exists
		if !utils.IsDryRun() {
			if err := utils.CreateDirectory("/opt/homebrew/etc/httpd/ssl"); err != nil {
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

		// Enable ssl_module in httpd.conf
		if !utils.IsDryRun() {
			if err := EnableSSLModuleInHttpdConf(); err != nil {
				return utils.LogError("Enabling SSL module in httpd.conf", err)
			}
		} else {
			utils.LogInfo("DRY RUN: Would enable SSL module in httpd.conf.")
		}

		// Restart Apache to apply changes
		if restartErr := RestartApache(); restartErr != nil {
			return fmt.Errorf("failed to restart Apache after enabling SSL module: %s", restartErr.Error())
		}

	} else {
		utils.LogSuccess("SSL certificates already exist.")
	}

	return nil
}

func EnableSSLModuleInHttpdConf() error {
	if utils.IsDryRun() {
		utils.LogInfo("DRY RUN: Would enable SSL module in Apache configuration.")
		return nil
	}

	utils.LogInfo("Enabling SSL module in Apache configuration...")

	file, err := os.Open(HttpdConfPath)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf: %s", err.Error())
	}
	defer file.Close()

	var lines []string
	sslModuleLoaded := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		// Check if the SSL module is already loaded
		if strings.Contains(line, "LoadModule ssl_module") && !strings.HasPrefix(line, "#") {
			sslModuleLoaded = true
		} else if strings.Contains(line, "LoadModule ssl_module") && strings.HasPrefix(line, "#") {
			// Uncomment the SSL module loading line if it's commented
			line = strings.TrimPrefix(line, "#")
			sslModuleLoaded = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read httpd.conf: %s", err.Error())
	}

	// Add the SSL module loading line if not found
	if !sslModuleLoaded {
		lines = append([]string{"LoadModule ssl_module lib/httpd/modules/mod_ssl.so"}, lines...)
		utils.LogSuccess("✔ Added SSL module loading line to httpd.conf.")
	}

	// Write the updated content back to httpd.conf
	file, err = os.OpenFile(HttpdConfPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf for writing: %s", err.Error())
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to httpd.conf: %s", err.Error())
		}
	}

	// uncomment the Include /opt/homebrew/etc/httpd/extra/httpd-ssl.conf line if it's commented
	includeSslConfAdded := false
	for _, line := range lines {
		if strings.Contains(line, "Include "+HttpdSSLConfPath) {
			includeSslConfAdded = true
			break
		}
	}
	if !includeSslConfAdded {
		lines = append(lines, "Include "+HttpdSSLConfPath)
		utils.LogSuccess("✔ Added Include " + HttpdSSLConfPath + " line to httpd.conf.")
	} else {
		utils.LogSuccess("Include " + HttpdSSLConfPath + " line already exists in httpd.conf.")
	}

	// Write the updated content back to httpd.conf
	file, err = os.OpenFile(HttpdConfPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf for writing: %s", err.Error())
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to httpd.conf: %s", err.Error())
		}
	}

	// read the httpd-ssl.conf file
	httpdSslConfFile, err := os.Open(HttpdSSLConfPath)
	if err != nil {
		return fmt.Errorf("failed to open httpd-ssl.conf: %s", err.Error())
	}
	defer httpdSslConfFile.Close()

	var httpdSslConfLines []string
	scanner = bufio.NewScanner(httpdSslConfFile)
	for scanner.Scan() {
		line := scanner.Text()
		httpdSslConfLines = append(httpdSslConfLines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read httpd-ssl.conf: %s", err.Error())
	}

	// add SSLCertificateFile and SSLCertificateKeyFile lines pointing to the SSL certificate and key files
	sslCertificateFileAdded := false
	for _, line := range lines {
		if strings.Contains(line, "SSLCertificateFile "+sslCertificateFile) {
			sslCertificateFileAdded = true
			break
		}
	}
	if !sslCertificateFileAdded {
		lines = append(lines, "SSLCertificateFile "+sslCertificateFile)
		utils.LogSuccess("✔ Added SSLCertificateFile line to httpd.conf.")
	} else {
		utils.LogSuccess("SSLCertificateFile line already exists in httpd.conf.")
	}

	sslCertificateKeyFileAdded := false
	for _, line := range lines {
		if strings.Contains(line, "SSLCertificateKeyFile "+sslCertificateKeyFile) {
			sslCertificateKeyFileAdded = true
			break
		}
	}
	if !sslCertificateKeyFileAdded {
		lines = append(lines, "SSLCertificateKeyFile "+sslCertificateKeyFile)
		utils.LogSuccess("✔ Added SSLCertificateKeyFile line to httpd.conf.")
	} else {
		utils.LogSuccess("SSLCertificateKeyFile line already exists in httpd.conf.")
	}

	// Write the updated content back to httpd.conf
	file, err = os.OpenFile(HttpdConfPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf for writing: %s", err.Error())
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to httpd.conf: %s", err.Error())
		}
	}

	utils.LogSuccess("SSL module enabled in httpd.conf.")
	return nil
}
