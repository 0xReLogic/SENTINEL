

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/0xReLogic/SENTINEL/config"
	"github.com/0xReLogic/SENTINEL/storage"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   cmdNameRun,
	Short: descRunShort,
	Long:  descRunLong,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		printBanner(cfg)

		// Initialize storage if configured
		var store storage.Storage
		if cfg.Storage.Type == "sqlite" && cfg.Storage.Path != "" {
			store, err = storage.NewSQLiteStorage(cfg.Storage.Path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to initialize storage: %v\n", err)
			} else {
				defer store.Close()
				fmt.Printf("Storage enabled: %s (retention: %d days)\n", cfg.Storage.Path, cfg.Storage.RetentionDays)
			}
		}

		stateManager := NewStateManager()

		workerCount := getWorkerCount()
		jobQueue := make(chan config.Service, workerCount)

		stop := make(chan os.Signal, 1)
		done := make(chan struct{})
		var workerWg sync.WaitGroup
		var schedulerWg sync.WaitGroup
		var mu sync.Mutex

		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Start cleanup goroutine if storage is enabled
		if store != nil && cfg.Storage.RetentionDays > 0 {
			go func() {
				ticker := time.NewTicker(24 * time.Hour)
				defer ticker.Stop()
				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						if err := store.Cleanup(cfg.Storage.RetentionDays); err != nil {
							fmt.Fprintf(os.Stderr, "Warning: Storage cleanup failed: %v\n", err)
						}
					}
				}
			}()
		}

		for i := 0; i < workerCount; i++ {
			workerWg.Add(1)
			go func() {
				defer workerWg.Done()
				for {
					select {
					case <-done:
						return
					case service, ok := <-jobQueue:
						if !ok {
							return
						}
						status := checker.CheckService(service.Name, service.URL, service.Timeout)
						
						mu.Lock()
						fmt.Println(status)
						mu.Unlock()

						// Save to storage if configured
						if store != nil {
							if err := store.SaveCheck(status); err != nil {
								fmt.Fprintf(os.Stderr, "Warning: Failed to save check to storage: %v\n", err)
							}
						}

						processNotifications(cfg, stateManager, status, service)
					}
				}
			}()
		}

		for _, service := range cfg.Services {
			jobQueue <- service
		}

		for _, service := range cfg.Services {
			service := service
			schedulerWg.Add(1)
			go func() {
				defer schedulerWg.Done()
				ticker := time.NewTicker(service.Interval)
				defer ticker.Stop()
				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						select {
						case jobQueue <- service:
						case <-done:
							return
						}
					}
				}
			}()
		}

		<-stop
		close(done)
		signal.Stop(stop)
		schedulerWg.Wait()
		close(jobQueue)
		workerWg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func printBanner(cfg *config.Config) {
	fmt.Println(bannerTitle)
	fmt.Printf(fmtLoadedServices, len(cfg.Services))
	fmt.Println(bannerExitInstruction)
	fmt.Println(separator)
}

func getWorkerCount() int {
	value, ok := os.LookupEnv(envWorkerCount)
	if !ok {
		return defaultWorkerCount
	}
	count, err := strconv.Atoi(value)
	if err == nil && count > 0 {
		return count
	}
	fmt.Fprintf(os.Stderr, msgInvalidWorkerCountEnv, envWorkerCount, value, defaultWorkerCount)
	return defaultWorkerCount
}