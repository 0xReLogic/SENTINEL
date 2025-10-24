// SENTINEL - A simple and effective monitoring system written in Go
// Repository: https://github.com/0xReLogic/SENTINEL

package main

import (
	"github.com/0xReLogic/SENTINEL/cmd"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("INFO: Could not load .env file (is that expected?): %v", err)
	}
	cmd.Execute()
}
