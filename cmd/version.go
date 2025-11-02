package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	versionNumber = "dev"
	commitHash    = "none"
	buildDate     = "unknown"
)

// SetVersionInfo sets the version information from main package
func SetVersionInfo(version, commit, date string) {
	versionNumber = version
	commitHash = commit
	buildDate = date
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  "Display version number, build date, commit hash, and Go version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SENTINEL version %s\n", versionNumber)
		fmt.Printf("Commit: %s\n", commitHash)
		fmt.Printf("Built: %s\n", buildDate)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
