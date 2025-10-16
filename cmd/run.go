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
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   cmdNameRun,
	Short: descRunShort,
	Long:  descRunLong,
	Run: func(cmd *cobra.Command, args []string) {
		// load configuration
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		// print startup banner
		printBanner(cfg)

		workerCount := getWorkerCount()
		jobQueue := make(chan config.Service, workerCount)

		// setup graceful shutdown handling
		stop := make(chan os.Signal, 1)
		done := make(chan struct{})
		var workerWg sync.WaitGroup
		var schedulerWg sync.WaitGroup
		var mu sync.Mutex

		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// start worker goroutines
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
						fmt.Printf(fmtTimestamp, time.Now().Format(timestampFormat),
							fmt.Sprintf(msgRunningServiceChecks, service.Name))
						fmt.Println(status)
						fmt.Println(separator)
						mu.Unlock()
					}
				}
			}()
		}

		// run the first check immediately for all services
		for _, service := range cfg.Services {
			jobQueue <- service
		}

		// schedule service checks
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

// printBanner displays the startup banner
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
