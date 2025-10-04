package cmd

import "time"

// application constants
const (
	// application name and metadata
	appName       = "sentinel"
	appRepository = "https://github.com/0xReLogic/SENTINEL"

	// command names
	cmdNameRun      = "run"
	cmdNameOnce     = "once"
	cmdNameValidate = "validate"

	// flag names
	flagConfig      = "config"
	flagConfigShort = "c"

	// default configuration
	defaultConfigFile = "sentinel.yaml"
	checkInterval     = 1 * time.Minute

	// timestamp format (Go reference time)
	timestampFormat = "2006-01-02 15:04:05"

	// allowed URL schemes
	schemeHTTP  = "http"
	schemeHTTPS = "https"

	// exit codes
	exitSuccess     = 0
	exitError       = 1
	exitConfigError = 2

	// display formatting
	separator  = "-----------------------------------"
	indent     = "  "
	listPrefix = "  - "

	// banner messages
	bannerTitle           = "SENTINEL Monitoring System"
	bannerExitInstruction = "Press Ctrl+C to exit"

	// validation messages
	msgValidationFailed   = "Configuration validation failed:\n"
	msgValidationSuccess  = "Configuration is valid"
	msgNoServicesDefined  = "No services defined"
	msgServicesConfigured = "Services configured:"

	// check messages
	msgRunningChecks = "Running service checks..."

	// error messages
	errLoadingConfig     = "Error loading configuration: %v\n"
	errInvalidConfigPath = "invalid config path: %w"
	errConfigNotFound    = "config file not found: %s\nCreate a %s file or use --%s flag"
	errServiceNameReq    = "service #%d: name is required"
	errServiceURLReq     = "service #%d (%s): URL is required"
	errServiceURLInvalid = "service #%d (%s): invalid URL format '%s'"

	// command descriptions
	descShort      = "A simple and effective monitoring system"
	descLong       = "SENTINEL monitors web services via HTTP and reports their status.\nPerfect for personal use or small teams needing lightweight monitoring.\n\nRepository: %s"
	descRunShort   = "Run continuous monitoring"
	descRunLong    = "Start SENTINEL in continuous monitoring mode. Checks run every minute."
	descOnceShort  = "Run checks once and exit"
	descOnceLong   = "Run service checks once and exit. Useful for cron jobs or CI/CD pipelines.\n\nExit codes:\n  %d - All services are UP\n  %d - One or more services are DOWN\n  %d - Configuration error"
	descValidShort = "Validate configuration file"
	descValidLong  = "Validate the configuration file for syntax and content errors."
	descConfigFlag = "path to configuration file"

	// message formats
	fmtLoadedServices           = "Loaded %d services to monitor\n"
	fmtLoadedServicesValidation = "Loaded %d services\n\n"
	fmtServiceListItem          = "  %d. %s - %s\n"
	fmtTimestamp                = "\n[%s] %s\n"
)
