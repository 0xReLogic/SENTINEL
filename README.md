# SENTINEL

<div align="center">
  <p><strong>SENTINEL</strong> - A simple and effective monitoring system written in Go.</p>
</div>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xReLogic/SENTINEL)](https://goreportcard.com/report/github.com/0xReLogic/SENTINEL)
[![Go Version](https://img.shields.io/github/go-mod/go-version/0xReLogic/SENTINEL)](https://github.com/0xReLogic/SENTINEL)
[![Build Status](https://img.shields.io/github/actions/workflow/status/0xReLogic/SENTINEL/go.yml?branch=main)](https://github.com/0xReLogic/SENTINEL/actions)
[![Release](https://img.shields.io/github/v/release/0xReLogic/SENTINEL)](https://github.com/0xReLogic/SENTINEL/releases)
[![Downloads](https://img.shields.io/github/downloads/0xReLogic/SENTINEL/total)](https://github.com/0xReLogic/SENTINEL/releases)

## Description

SENTINEL is a simple monitoring system written in Go. This application can monitor the status of various web services via HTTP and report their status periodically. It's ideal for personal use or small teams that need a lightweight and easily configurable monitoring solution.

## Features

- Monitor various web services via HTTP
- Simple configuration using YAML format
- UP/DOWN status reporting with response time
- Automatic checks every minute
- Concurrency for efficient checking
- Flexible CLI with various commands

## Installation

### Prerequisites

- Go 1.21 or newer

### From Source

```bash
# Clone repository
git clone https://github.com/0xReLogic/SENTINEL.git
cd SENTINEL

# Build application
make build

# Or use Go directly
go build -o sentinel
```

### Using Go Install

```bash
go install github.com/0xReLogic/SENTINEL@latest
```

## How to Use

1. Create a `sentinel.yaml` configuration file (see example below)
2. Run SENTINEL with one of the following commands:

```bash
# Run continuous checks
./sentinel run

# Run a single check
./sentinel once

# Validate configuration file
./sentinel validate

# Display help
./sentinel --help

# Use custom configuration file
./sentinel run --config /path/to/config.yaml
```

## Configuration Structure

The `sentinel.yaml` configuration file has the following format:

```yaml
# SENTINEL Configuration File
services:
  - name: "Google"
    url: "https://www.google.com"
  - name: "GitHub"
    url: "https://github.com"
  - name: "Example"
    url: "https://example.com"
```

## Project Structure

```
SENTINEL/
├── checker/       # Package for service checking
├── cmd/           # CLI commands
├── config/        # Package for configuration management
├── notifier/      # Notification for Telegram
├── main.go        # Main program file
├── Makefile       # Makefile for easier build and test
├── go.mod         # Go module definition
├── go.sum         # Dependencies checksum
├── sentinel.yaml  # Example configuration file
├── LICENSE        # MIT License
├── README.md      # Main documentation
└── CONTRIBUTING.md # Contribution guidelines
```

## Roadmap

- [x] Phase 1: Core Engine - Basic implementation
  - [x] 1.1: Create basic service checking functionality
  - [x] 1.2: Implement configuration loading from YAML
  - [x] 1.3: Build simple CLI interface

- [x] Phase 2: Making it "Smart" - Better concurrency and CLI
  - [x] 2.1: Add concurrent service checking
  - [x] 2.2: Implement more robust CLI with subcommands
  - [x] 2.3: Add proper error handling and logging

- [x] Phase 3: Advanced Features - Notifications and advanced check types
  - [x] 3.1: Add notification system (email, webhook)
  - [x] 3.2: Implement advanced check types (TCP, ICMP)
  - [x] 3.3: Add custom thresholds and alerting rules

- [x] Phase 4: Open Source Readiness - Documentation and testing
  - [x] 4.1: Complete code documentation and examples
  - [x] 4.2: Add comprehensive test suite with high coverage
  - [x] 4.3: Create contribution guidelines and code of conduct
  - [x] 4.4: Add CI/CD pipeline with GitHub Actions

- [ ] Phase 5: Enterprise Features - Advanced monitoring capabilities
  - [ ] 5.1: Implement distributed monitoring with agent architecture
  - [ ] 5.2: Add support for custom plugins and extensions
  - [ ] 5.3: Create dashboard for real-time monitoring visualization
  - [ ] 5.4: Implement historical data storage and analysis

- [ ] Phase 6: Integration and Ecosystem - Connect with other tools
  - [ ] 6.1: Add integration with popular alerting systems (PagerDuty, OpsGenie)
  - [ ] 6.2: Implement support for metrics export to Prometheus/Grafana
  - [ ] 6.3: Create API for third-party integrations
  - [ ] 6.4: Develop SDKs for popular programming languages

- [ ] Phase 7: Scalability and Performance - Enterprise-grade monitoring
  - [ ] 7.1: Optimize for high-volume monitoring environments
  - [ ] 7.2: Implement clustering for horizontal scaling
  - [ ] 7.3: Add support for multi-region monitoring
  - [ ] 7.4: Create enterprise deployment guides and best practices

## Contribution

Contributions are greatly appreciated! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on the pull request submission process.

## Releasing New Versions

To release a new version:

1. Create a new tag with semantic versioning:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

2. The GitHub Actions workflow will automatically:
   - Build binaries for multiple platforms (Linux, Windows, macOS)
   - Create a GitHub release with the binaries attached
   - Generate release notes based on commit messages

Binary releases will be available at: https://github.com/0xReLogic/SENTINEL/releases

## Local Release

You can create a local release using the provided scripts:

### Windows
```bash
# Run the build-release.bat script
.\build-release.bat
```

### Linux/macOS
```bash
# Give execute permission
chmod +x build-release.sh

# Run the script
./build-release.sh
```
## Notifications

SENTINEL can send real-time alerts when a service goes down or recovers.

### Telegram Setup

To receive notifications in a Telegram chat, follow these steps:

**Step 1: Get Telegram Credentials**

1.  **Bot Token**:
    * Open Telegram and start a chat with the official **`@BotFather`**.
    * Send the `/newbot` command and follow the prompts to create your bot.
    * `@BotFather` will give you a unique **Bot Token**. This is a secret, so don't share it publicly.

2.  **Chat ID**:
    * Start a chat with your newly created bot by finding it and sending the `/start` command. **This is a required step.**
    * Next, start a chat with the bot **`@userinfobot`**.
    * Send it a message, and it will reply with your user information, including your **Chat ID**.

**Step 2: Configure SENTINEL**

1.  **Create a `.env` file** in the root directory of the project. This file will securely store your secrets. **Important**: Add `.env` to your `.gitignore` file to avoid committing secrets to your repository.

    ```
    # .env
    # Variables for Telegram notifications
    TELEGRAM_BOT_TOKEN="123456:ABC-DEF1234ghIkl-JAS0987-aB5-qwer"
    TELEGRAM_CHAT_ID="123456789"
    ```

2.  **Update your `sentinel.yaml`** to enable and configure Telegram notifications. The `${VAR_NAME}` syntax will securely load the values from your `.env` file.

    ```yaml
    # sentinel.yaml
    
    notifications:
      telegram:
        enabled: true
        bot_token: "${TELEGRAM_BOT_TOKEN}"
        chat_id: "${TELEGRAM_CHAT_ID}"
        notify_on:
          - down
          - recovery
    ```

#### Example Notification Messages

When a service status changes, you will receive a message in your configured chat.

**Service Down**
> 🔴 **Service DOWN**
> **Name:** My Failing API
> **URL:** https://api.example.com/health
> **Error:** connection timeout
> **Time:** 2025-10-12 10:10:00

**Service Recovered**
> 🟢 **Service RECOVERED**
> **Name:** My Failing API
> **URL:** https://api.example.com/health
> **Downtime:** 5m 30s
> **Time:** 2025-10-12 10:15:30

The script will automatically build binaries for all platforms (Linux, Windows, macOS) and place them in the `./dist` folder.
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors who have helped with this project
- Inspiration from various monitoring systems such as Prometheus, Nagios, and Uptime Robot
