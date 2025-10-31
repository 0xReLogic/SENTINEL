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
- Customizable check intervals and timeouts per service
- UP/DOWN status reporting with response time
- Automatic checks on configurable intervals
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

## Docker Deployment

### Using Docker Compose (Recommended)
```bash
# Copy environment variables
cp .env.example .env
# Edit .env with your tokens - get these from mentioned steps 

# Start SENTINEL
docker compose up -d

# View logs
docker compose logs -f

# Stop SENTINEL
docker compose down
```

### Using Docker Directly
```bash
docker build -t sentinel .

# Run container
docker run --rm --env-file .env -v ./sentinel.yaml:/app/sentinel.yaml -p 8080:8080 sentinel
```

### Configuration
Mount your `sentinel.yaml` config file as a volume to customize which services to monitor.

## Configuration Structure

The `sentinel.yaml` configuration file has the following format:

```yaml
# SENTINEL Configuration File
services:
  - name: "Google"
    url: "https://www.google.com"
    interval: 30s       # optional, default is 1m
    timeout: 3s         # optional, default is 5s
  - name: "GitHub"
    url: "https://github.com"
    interval: 2m
  - name: "Example"
    url: "https://example.com"
    # No interval/timeout defined -> defaults apply
```

If `interval` or `timeout` are omitted, SENTINEL falls back to the defaults of `1m`
and `5s` respectively.

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

### Discord Setup

To receive notifications in a Discord channel via webhook:

**Step 1: Get Discord Webhook URL**

1. Open your Discord server and go to **Server Settings** → **Integrations** → **Webhooks**
2. Click **"New Webhook"** or **"Create Webhook"**
3. Give your webhook a name (e.g., "SENTINEL Monitor")
4. Select the channel where you want to receive notifications
5. Click **"Copy Webhook URL"** - this is your webhook URL

**Step 2: Configure SENTINEL**

1. **Add to your `.env` file**:

    ```
    # .env
    # Variables for Discord notifications
    DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/123456789/abcdefghijklmnop"
    ```

2. **Update your `sentinel.yaml`**:

    ```yaml
    # sentinel.yaml
    
    notifications:
      telegram:
        enabled: false  # Can enable both at the same time
        bot_token: "${TELEGRAM_BOT_TOKEN}"
        chat_id: "${TELEGRAM_CHAT_ID}"
        notify_on:
          - down
          - recovery
      discord:
        enabled: true
        webhook_url: "${DISCORD_WEBHOOK_URL}"
        notify_on:
          - down
          - recovery
    ```

#### Example Discord Notification Messages

Discord notifications use rich embeds with color coding:

**Service Down** (Red embed)
- **Title:** 🔴 Service DOWN
- **Fields:**
  - Service: My Failing API
  - URL: https://api.example.com/health
  - Error: connection timeout
- **Timestamp:** 2025-10-12T10:10:00Z

**Service Recovered** (Green embed)
- **Title:** 🟢 Service RECOVERED
- **Fields:**
  - Service: My Failing API
  - URL: https://api.example.com/health
  - Downtime: 5m 30s
- **Timestamp:** 2025-10-12T10:15:30Z

**Note:** You can enable both Telegram and Discord notifications simultaneously. SENTINEL will send alerts to all enabled notification channels.

The script will automatically build binaries for all platforms (Linux, Windows, macOS) and place them in the `./dist` folder.
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributors

Thanks to all the amazing people who have contributed to SENTINEL. 🎉

<a href="https://github.com/0xReLogic/SENTINEL/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=0xReLogic/SENTINEL" />
</a>


## Acknowledgments

- Inspiration from various monitoring systems such as Prometheus, Nagios, and Uptime Robot
