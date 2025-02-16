This project is my personal endeavor to learn and explore the Go programming language (Golang). While the tool is functional and built with care, it primarily serves as a way for me to deepen my understanding of various Go concepts and build practical skills. Read the full story in the [ABOUT](ABOUT.md) section.

**Note:** This tool is named `localhost`, but it is not related to the `localhost` networking hostname (127.0.0.1). Instead, this tool is designed to help developers manage local domains and virtual host configurations on macOS.


# LocalHost - multiple local domains on MacOS

[![Go Report Card](https://goreportcard.com/badge/github.com/liviu-hariton/localhost)](https://goreportcard.com/report/github.com/liviu-hariton/localhost)
![Go Version](https://img.shields.io/github/go-mod/go-version/liviu-hariton/localhost)
![License](https://img.shields.io/github/license/liviu-hariton/localhost)
![Last Commit](https://img.shields.io/github/last-commit/liviu-hariton/localhost)
![GitHub Release](https://img.shields.io/github/v/release/liviu-hariton/localhost)
![GitHub Downloads](https://img.shields.io/github/downloads/liviu-hariton/localhost/total)
![Platform](https://img.shields.io/badge/platform-macOS-blue)

LocalHost is a utility for web developers to set up and manage multiple local domains on macOS. It's meant for those users who want to relly on vanilla tools and avoid using memory hungry tools like Docker. It automates tedious tasks like:
* creating and loading virtual hosts files
* modifying `/etc/hosts` file
* setting up the MySQL database server and / or the PHP interpreter
* setting up SSL for local domains

### Key Features
* automatic installation of MySQL database server and, also, for the PHP interpreter
    * the lastes versions are automatically installed and launched by using [Homebrew](https://brew.sh/)
* create and delete local domains (with interactive prompts to avoid accidental changes)
* automatic configuration of `/etc/hosts` and virtual hosts
* SSL setup with self-signed certificates
* built-in dry-run mode for safe experimentation, whitout actually performing any action

### Ideal For
* web developers who often work with multiple projects requiring custom local domains (e.g., myproject.local).
* freelancers and solo developers who often juggle multiple projects for different clients, each with its own domain and directory setup
* beginner developers who might not yet be familiar with configuring Apache, editing /etc/hosts, or troubleshooting local domain setups

## Table of Contents

* [Requirements](#requirements)
* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Usage](#usage)
    * [See the Help documentation](#see-the-help-documentation)
    * [Create a local domain](#create-a-local-domain)
    * [List available local domains](#list-available-local-domains)
    * [Remove an existing local domain](#remove-an-existing-local-domain)
    * [Dry-Run mode](#dry-run-mode)
* [Uninstallation](#uninstallation)
* [Build it yourself](#build-it-yourself)
* [License](#license)
* [Disclaimer](#disclaimer)

### Requirements

* a macOS operating system
    * this tool was tested on the latest macOS Sequoia version 15.3 but, in theory, it should work on previous versions
    * the binary works on both Intel (`amd64`) and Apple Silicon (`arm64`):
* make sure you have [Homebrew](https://brew.sh/) up and running on your system

### Prerequisites

* macOS includes a pre-installed version of Apache. However, you may check to see if the `/usr/local/etc/httpd/extra/` directory exists as expected
    * here, the folder `vhosts` will be created for keeping the virtual host configuration files that will be created
* the standard Apache configuration is located at `/usr/local/etc/httpd/httpd.conf`
* the self-signed certificates will be stored in
    * `/etc/apache2/ssl/server.crt`
    * `/etc/apache2/ssl/server.key`

### Installation

**Manual installation**

1. Go to the [Releases](https://github.com/liviu-hariton/localhost/releases) page and download the latest version.

```bash
curl -LO https://github.com/liviu-hariton/localhost/releases/download/1.0.0/localhost
```

2. Make the binary executable

```bash
chmod +x ./localhost
```

3. Move the binary to `/usr/local/bin`

```bash
sudo mv localhost /usr/local/bin/
```

4. Verify the installation

```bash
localhost --version
```

### See the Help documentation

List all available commands

```bash
localhost help
```

### Create a local domain

You nee to provide two parameters to the command:
* `-domain` - this is your new local domain name - the one that you'll access in your browser (e.g., `myproject.local`)
* `-doc_root` - this it the path on disk were your project's files will reside (e.g., `/path/on/disk/to/your/project`)
    * if the path does not exists, it will be created automatically for you

You can, also, add the `--no-dns-reset` flag to skip the local DNS cache flushing and resetting the `mDNSResponder`

```bash
localhost create -domain=myproject.local -doc_root=/path/to/myproject
```

#### How it works

* checks if Apache is installed and, if Apache is not running, attempts to restart it
* checks if MySQL is installed and, if not, tries to install it via Homebrew. Next, it attempts to start it as a background service
    * it will install [the latest MySQL version available in Homebrew](https://formulae.brew.sh/formula/mysql#default)
* checks if PHP is installed and, if not, tries to install it via Homebrew. Next, it attempts to start it as a background service
    * it will install [the latest PHP version available in Homebrew](https://formulae.brew.sh/formula/php#default)
    * also, it will enable the PHP module in the Apache's standard configuration
* adds the new required entry into the `/etc/hosts` file
* checks if virtual hosts are enabled in your Apache configuration and, if so, creates the new virtual host configuration
* ensures the SSL certificate and key files exist, generating them if necessary (self-signed)

### List available local domains

In order to see what local domains (configurations) are available at any time, run:

```bash
localhost list
```

You will get an output like this:

```bash
Configured domains:
[DEBUG] myproject.local.conf
[DEBUG] someotherproject.local.conf
```

### Remove an existing local domain

In order to remove an existing local domain, run:
* you can, also, add the `--no-dns-reset` flag to skip the local DNS cache flushing and resetting the `mDNSResponder`

```bash
localhost delete -domain=myproject.local
```

#### How it works

* removes the local domain entry from `/etc/hosts`
* deletes the corresponding virtual host configuration file, previously created
* restarts Apache and flushes the DNS cache (if the `--no-dns-reset` flag is not set)

### Dry-Run mode

You can simulates actions without making any actual changes to your system (so that you can preview the actions about to be applied) by adding the `--dry-run` flag

```bash
localhost create -domain=myproject.local -doc_root=/path/to/myproject --dry-run
```

You will get an output like this:

```bash
[INFO] Dry Run mode detected. Skipping privilege escalation.
[INFO] Running in Dry Run mode: No changes will be made.
[INFO] Starting setup for domain: myproject.local

Starting system checks...
[INFO] Checking Apache setup...
✔ Apache is installed.
✔ Apache is running.
[INFO] Checking MySQL setup...
✔ MySQL is installed.
✔ MySQL is running.
[INFO] Checking PHP setup...
[SUCCESS] PHP is installed.
[SUCCESS] PHP is working correctly.
[SUCCESS] All checks passed successfully!
Modifying hosts file for domain: myproject.local
DRY RUN: Would add the domain to the hosts file.
Ensuring vhosts are enabled and adding virtual host...
DRY RUN: Would update httpd.conf with new content.
✔ Virtual hosts wildcard line already exists in httpd.conf.
DRY RUN: Would create directory: /usr/local/etc/httpd/extra/vhosts/
DRY RUN: Would create directory: /path/to/myproject/_logs/myproject.local/ssl
DRY RUN: Would create directory: /path/to/myproject/public
DRY RUN: Would write the dummy index.php file.
DRY RUN: Would write the virtual host configuration file.
[INFO] Restarting Apache to apply changes...
DRY RUN: Would restart Apache server and flush the DNS cache.
[INFO] Checking for SSL certificates...
[SUCCESS] SSL certificates already exist.
[SUCCESS] All changes applied successfully!
[INFO] You should now be able to access your new project at http://myproject.local or https://myproject.local
```

### Uninstallation

You can uninstall the **LocalHost** utility from your system by following these steps:

1. Remove the binary from `/usr/local/bin`:

```bash
sudo rm /usr/local/bin/localhost
```

2. Clean up any leftover configurations

To remove all configurations created by the tool:

```bash
sudo rm -rf /usr/local/etc/httpd/extra/vhosts/*.conf
```
Alternatively, you can delete specific virtual host files:

```bash
sudo rm /usr/local/etc/httpd/extra/vhosts/myproject.local.conf
```

3. Remove entries from `/etc/hosts`

```bash
sudo nano /etc/hosts
```
Find and delete lines that reference the local domains (e.g., `127.0.0.1 myproject.local`). Save the file and exit (`:wq`) .

## Build it yourself

This guide will walk you through the steps to clone, build, and run the localhost utility from source. You can customize or contribute to the tool by following these instructions.

Before building the tool, ensure you have the following installed:
* Go (Golang)
    * download and install the latest version of Go from the [official website](https://go.dev/)
    * verify the installation with `go version`
* Git
    * install Git from [git-scm.com](git-scm.com)
    * verify the installation with `git --version`
* Xcode command line tools (macOS-specific)
    * run this in your terminal `xcode-select --install`

Next, clone the official repository to your local machine

```bash
git clone https://github.com/liviu-hariton/localhost.git
```

and navigate into the project directory

```bash
cd localhost
```

Now, you have two options:
* build a native binary, for your current macOS architecture
* build an universal binary that runs natively on both Apple Silicon and Intel-based Macs

### Build a native binary

```bash
go build -o localhost
```

This will create a binary named `localhost` in the current directory. You can run the tool immediately:

```bash
./localhost --version
```

### Build an universal binary (for both Apple Silicon and Intel)

If you want a single binary that runs natively on both Apple Silicon and Intel-based Macs, you need to build architecture-specific binaries and merge them using `lipo`

First, build an ARM64 (Apple Silicon) binary:

```bash
GOOS=darwin GOARCH=arm64 go build -o localhost-arm64
```

Then, build an AMD64 (Intel) binary:

```bash
GOOS=darwin GOARCH=amd64 go build -o localhost-amd64
```

And, finally, merge them together into an universal binary:

```bash
lipo -create -output localhost localhost-arm64 localhost-amd64
```

You can check the architectures included in the universal binary you've just created by running:

```bash
lipo -info localhost
```

At this point, no matter the binary version chosen, you can make the `localhost` tool available globally, by moving it to a directory in your PATH, such as `/usr/local/bin`

```bash
sudo mv localhost /usr/local/bin/
```

Now you can test it by running:

```bash
localhost --version
```

### License
This library is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for more details.

### Disclaimer
This tool, **LocalHost**, is provided "as is" and is intended for use by developers on macOS for managing local development environments. See the [DISCLAIMER](DISCLAIMER.md) file for more details.
