package metrics

import (
	"testing"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

const (
	testServiceName = "Test Service"
	testServiceURL  = "https://example.com"
)

func TestRecordCheckSuccess(t *testing.T) {
	status := checker.ServiceStatus{
		Name:         testServiceName,
		URL:          testServiceURL,
		IsUp:         true,
		ResponseTime: 100 * time.Millisecond,
		StatusCode:   200,
	}

	RecordCheck(status)

	// Verify service_up gauge is 1
	value := testutil.ToFloat64(ServiceUp.WithLabelValues(testServiceName, testServiceURL))
	if value != 1.0 {
		t.Errorf("Expected service_up to be 1, got %f", value)
	}
}

func TestRecordCheckFailure(t *testing.T) {
	status := checker.ServiceStatus{
		Name:         "Failed Service",
		URL:          "https://failed.example.com",
		IsUp:         false,
		ResponseTime: 5 * time.Second,
		StatusCode:   500,
	}

	RecordCheck(status)

	// Verify service_up gauge is 0
	value := testutil.ToFloat64(ServiceUp.WithLabelValues("Failed Service", "https://failed.example.com"))
	if value != 0.0 {
		t.Errorf("Expected service_up to be 0, got %f", value)
	}
}

func TestStatusCodeToString(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{200, "2xx"},
		{201, "2xx"},
		{301, "3xx"},
		{404, "4xx"},
		{500, "5xx"},
		{0, "unknown"},
	}

	for _, tt := range tests {
		result := statusCodeToString(tt.code)
		if result != tt.expected {
			t.Errorf("statusCodeToString(%d) = %s, expected %s", tt.code, result, tt.expected)
		}
	}
}
