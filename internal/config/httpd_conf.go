package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HttpdConfPath defines the path to the Apache main configuration file.
const HttpdConfPath = "/usr/local/etc/httpd/httpd.conf"

// EnsureVhostsEnabled ensures that the httpd.conf file includes the vhosts file.
func EnsureVhostsEnabled() error {
	file, err := os.Open(HttpdConfPath)
	if err != nil {
		return fmt.Errorf("failed to open httpd.conf: %s", err.Error())
	}
	defer file.Close()

	var lines []string
	includeVhostsWildcard := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the vhosts wildcard line already exists
		if strings.Contains(line, "Include /usr/local/etc/httpd/extra/vhosts/*.conf") {
			includeVhostsWildcard = true
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read httpd.conf: %s", err.Error())
	}

	// If the wildcard line doesn't exist, add it after the default vhosts line
	if !includeVhostsWildcard {
		fmt.Println("✔ Adding 'Include /usr/local/etc/httpd/extra/vhosts/*.conf' to httpd.conf")
		lines = append(lines, "Include /usr/local/etc/httpd/extra/vhosts/*.conf")
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

	if !includeVhostsWildcard {
		fmt.Println("✔ Virtual hosts wildcard line added to httpd.conf.")
	} else {
		fmt.Println("✔ Virtual hosts wildcard line already exists in httpd.conf.")
	}
	return nil
}

// AddVirtualHost creates a new virtual host configuration for the domain.
func AddVirtualHost(domain, documentRoot string) error {
	// Define the path for the new vhost config file
	vhostsDir := "/usr/local/etc/httpd/extra/vhosts/"
	vhostFile := vhostsDir + domain + ".conf"

	// Derive log paths based on the document root
	baseLogDir := fmt.Sprintf("%s/_logs/%s", documentRoot, domain)
	errorLogDir := fmt.Sprintf("%s/error_log", baseLogDir)
	accessLogDir := fmt.Sprintf("%s/access_log", baseLogDir)
	sslErrorLogDir := fmt.Sprintf("%s/ssl/error_log", baseLogDir)
	sslAccessLogDir := fmt.Sprintf("%s/ssl/access_log", baseLogDir)

	// Ensure the vhosts directory exists
	if err := os.MkdirAll(vhostsDir, 0755); err != nil {
		return fmt.Errorf("failed to create vhosts directory: %s", err.Error())
	}

	// Ensure the log directories exist
	if err := os.MkdirAll(fmt.Sprintf("%s/ssl", baseLogDir), 0755); err != nil {
		return fmt.Errorf("failed to create log directories: %s", err.Error())
	}

	// Ensure the public directory exists
	publicDir := fmt.Sprintf("%s/public", documentRoot)
	if err := os.MkdirAll(publicDir, 0755); err != nil {
		return fmt.Errorf("failed to create public directory: %s", err.Error())
	}

	// Write the dummy index.php file
	indexPhpFile := fmt.Sprintf("%s/index.php", publicDir)
	indexPhpContent := fmt.Sprintf("<?php\necho 'It worked! You are on %s domain.';\n", domain)

	if err := os.WriteFile(indexPhpFile, []byte(indexPhpContent), 0644); err != nil {
		return fmt.Errorf("failed to write to index.php file '%s': %s", indexPhpFile, err.Error())
	}
	fmt.Printf("✔ Dummy index.php file created at '%s'.\n", indexPhpFile)

	// Open the file for writing
	file, err := os.OpenFile(vhostFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create vhost file '%s': %s", vhostFile, err.Error())
	}
	defer file.Close()

	// Virtual host configuration template
	vhostConfig := fmt.Sprintf(`
<VirtualHost %s:80>
    ServerName %s
    DocumentRoot "%s/public"
    ErrorLog "%s"
    CustomLog "%s" common

    <Directory "%s">
        Options FollowSymLinks Multiviews Indexes
        MultiviewsMatch Any
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>

<VirtualHost %s:443>
    ServerName %s
    DocumentRoot "%s/public"
    SSLEngine on
    SSLCipherSuite ALL:!ADH:!EXPORT56:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP:+eNULL
    SSLCertificateFile /etc/apache2/ssl/server.crt
    SSLCertificateKeyFile /etc/apache2/ssl/server.key
    ErrorLog "%s"
    CustomLog "%s" common

    <Directory "%s">
        Options FollowSymLinks Multiviews Indexes
        MultiviewsMatch Any
        AllowOverride All
        Require all granted
    </Directory>
</VirtualHost>
`, domain, domain, documentRoot, errorLogDir, accessLogDir, documentRoot, domain, domain, documentRoot, sslErrorLogDir, sslAccessLogDir, documentRoot)

	// Write the configuration to the file
	if _, err := file.WriteString(vhostConfig); err != nil {
		return fmt.Errorf("failed to write to vhost file '%s': %s", vhostFile, err.Error())
	}

	fmt.Printf("✔ Virtual host configuration for '%s' created at '%s'.\n", domain, vhostFile)
	return nil
}
