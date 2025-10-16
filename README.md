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
