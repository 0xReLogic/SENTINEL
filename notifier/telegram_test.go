package notifier

import (
	"testing"
	"time"
)

// TestFormatDown verifies that the "service down" message is formatted correctly.
func TestFormatDown(t *testing.T) {
	name := "Test Down Message"
	url := "https://api.test-service.com"
	errorMsg := "connection timeout (error-code: 123)"
	checkTime := time.Date(2025, 10, 11, 22, 30, 0, 0, time.UTC)
	expected := `🔴 *Service DOWN*
*Name:* Test Down Message
*URL:* https://api\.test\-service\.com
*Error:* connection timeout \(error\-code: 123\)
*Time:* 2025\-10\-11 22:30:00`
	actual := FormatDownMessage(name, url, errorMsg, checkTime)

	if actual != expected {
		t.Errorf("FormatDownMessage() failed:\nExpected:\n%s\n\nGot:\n%s", expected, actual)
	}
}

// TestFormatRecovery verifies that the "service recovered" message is formatted correctly.
func TestFormatRecovery(t *testing.T) {
	name := "Test Recovery Message"
	url := "https://api.test-service.com"
	downtime := 5*time.Minute + 30*time.Second
	checkTime := time.Date(2025, 10, 11, 22, 30, 0, 0, time.UTC)
	expected := `🟢 *Service RECOVERED*
*Name:* Test Recovery Message
*URL:* https://api\.test\-service\.com
*Downtime:* 5m30s
*Time:* 2025\-10\-11 22:30:00`

	actual := FormatRecoveryMessage(name, url, downtime, checkTime)
	if actual != expected {
		t.Errorf("FormatRecoveryMessage() failed:\nExpected:\n%s\nGot:\n%s", expected, actual)
	}
}
