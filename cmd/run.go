

package cmd

import (
	"fmt"
	"log"
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
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, errLoadingConfig, err)
			os.Exit(exitConfigError)
		}

		printBanner(cfg)

		stateManager := NewStateManager()

		workerCount := getWorkerCount()
		jobQueue := make(chan config.Service, workerCount)

		stop := make(chan os.Signal, 1)
		done := make(chan struct{})
		var workerWg sync.WaitGroup
		var schedulerWg sync.WaitGroup
		var mu sync.Mutex

		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

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

						// Process Telegram notifications
						if cfg.Notifications.Telegram.Enabled {
							action := stateManager.ProcessStatus(status, service, cfg.Notifications.Telegram)
							switch action.Action {
							case NotifyDown:
								log.Printf("INFO: Service '%s' is DOWN. Preparing Telegram notification.", status.Name)
								NotifyServiceDown(cfg.Notifications.Telegram, status, time.Now())
							case NotifyRecovery:
								log.Printf("INFO: Service '%s' has RECOVERED. Preparing Telegram notification.", status.Name)
								NotifyServiceRecovery(cfg.Notifications.Telegram, status, action.Downtime, time.Now())
							}
						}

						// Process Discord notifications
						if cfg.Notifications.Discord.Enabled {
							tempCfg := config.TelegramConfig{
								Enabled:  cfg.Notifications.Discord.Enabled,
								NotifyOn: cfg.Notifications.Discord.NotifyOn,
							}
							action := stateManager.ProcessStatus(status, service, tempCfg)
							switch action.Action {
							case NotifyDown:
								log.Printf("INFO: Service '%s' is DOWN. Preparing Discord notification.", status.Name)
								NotifyDiscordServiceDown(cfg.Notifications.Discord, status, time.Now())
							case NotifyRecovery:
								log.Printf("INFO: Service '%s' has RECOVERED. Preparing Discord notification.", status.Name)
								NotifyDiscordServiceRecovery(cfg.Notifications.Discord, status, action.Downtime, time.Now())
							}
						}
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