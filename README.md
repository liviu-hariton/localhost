# LocalHost - multiple local domains on MacOS

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
    * [See the Help documentation](#)
    * [Create a local domain](#)
    * [List available local domains](#)
    * [Remove an existing local domain](#)
    * [Dry-Run mode](#)
* [Uninstallation](#uninstallation)
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

### Installation

**Manual installation**

1. Go to the [Releases](https://github.com/liviu-hariton/localhost/releases) page and download the latest version.

```bash
curl -LO https://github.com/liviu-hariton/localhost/releases/download/v1.0.0/localhost.tar.gz
```

2. Extract the binary

```bash
tar -xvzf localhost.tar.gz
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

```bash
localhost create -domain=myproject.local -doc_root=/path/to/myproject
```

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

```bash
localhost delete -domain=myproject.local
```

### Dry-Run mode

You can simulates actions without making any actual changes to your system (so that you can preview the actions about to be applied) by adding the `--dry-run` flag

```bash
localhost create -domain=myproject.local -doc_root=/path/to/myproject --dry-run
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

### License
This library is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for more details.

### Disclaimer
This tool, **LocalHost**, is provided "as is" and is intended for use by developers on macOS for managing local development environments. See the [DISCLAIMER](DISCLAIMER.md) file for more details.
