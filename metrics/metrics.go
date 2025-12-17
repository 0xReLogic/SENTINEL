// Package metrics provides Prometheus metrics for SENTINEL
package metrics

import (
	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ServiceUp indicates if a service is up (1) or down (0)
	ServiceUp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "sentinel",
			Name:      "service_up",
			Help:      "Service is up (1) or down (0)",
		},
		[]string{"service", "url"},
	)

	// ResponseTime tracks HTTP response time in seconds
	ResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sentinel",
			Name:      "response_time_seconds",
			Help:      "HTTP response time in seconds",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "url"},
	)

	// ChecksTotal counts total number of checks performed
	ChecksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sentinel",
			Name:      "checks_total",
			Help:      "Total number of checks performed",
		},
		[]string{"service", "status"},
	)

	// HTTPStatusTotal counts HTTP status codes received
	HTTPStatusTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sentinel",
			Name:      "http_status_total",
			Help:      "HTTP status codes received",
		},
		[]string{"service", "code"},
	)
)

// RecordCheck updates all metrics based on a service check result
func RecordCheck(status checker.ServiceStatus) {
	// Update service up/down gauge
	upValue := 0.0
	if status.IsUp {
		upValue = 1.0
	}
	ServiceUp.WithLabelValues(status.Name, status.URL).Set(upValue)

	// Record response time
	ResponseTime.WithLabelValues(status.Name, status.URL).Observe(status.ResponseTime.Seconds())

	// Increment check counter
	checkStatus := "failure"
	if status.IsUp {
		checkStatus = "success"
	}
	ChecksTotal.WithLabelValues(status.Name, checkStatus).Inc()

	// Record HTTP status code (only if we got a response)
	if status.StatusCode > 0 {
		HTTPStatusTotal.WithLabelValues(status.Name, statusCodeToString(status.StatusCode)).Inc()
	}
}

func statusCodeToString(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "2xx"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	case code >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
