// SENTINEL - A simple and effective monitoring system written in Go
// Repository: https://github.com/0xReLogic/SENTINEL

package main

import (
	"log"

	"github.com/0xReLogic/SENTINEL/cmd"
	"github.com/joho/godotenv"
)

// Build information (set via ldflags during build)
var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	// Set version info for cmd package
	cmd.SetVersionInfo(version, commit, buildDate)

	err := godotenv.Load()
	if err != nil {
		log.Printf("INFO: Could not load .env file (is that expected?): %v", err)
	}
	cmd.Execute()
}
