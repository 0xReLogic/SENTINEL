// in tests/system_test.go

package system_test

import (
	"testing"
	"time"

	"github.com/0xReLogic/SENTINEL/checker"
	"github.com/0xReLogic/SENTINEL/cmd"
	"github.com/0xReLogic/SENTINEL/config"
)

const testCheckInterval = 1 * time.Minute

var mockTelegramCfg = config.TelegramConfig{
	NotifyOn: []string{"down", "recovery"},
}

func TestProcessStatus_Transitions_Integration(t *testing.T) {

	// 1. Setup the Manager from the CMD package
	sm := cmd.NewStateManager()
	serviceURL := "https://testservice.com/status"
	serviceName := "TestService"

	
	// We use the test constant for the interval to keep the test logic consistent.
	mockService := config.Service{
		Name:     serviceName,
		URL:      serviceURL,
		Interval: testCheckInterval,
	}


	downStatus := checker.ServiceStatus{Name: serviceName, URL: serviceURL, IsUp: false}
	upStatus := checker.ServiceStatus{Name: serviceName, URL: serviceURL, IsUp: true}

	// --- SCENARIO 1: Initial Check (UP) -> No Action ---
	t.Run("Initial_UP_NoAction", func(t *testing.T) {

		action := sm.ProcessStatus(upStatus, mockService, mockTelegramCfg)
		if action.Action != cmd.NoAction {
			t.Errorf("Expected initial UP check to be NoAction, got %v", action.Action)
		}
	})

	// --- SCENARIO 2: UP -> DOWN Transition (NotifyDown) ---
	t.Run("UP_to_DOWN_NotifyDown", func(t *testing.T) {
		
		action := sm.ProcessStatus(downStatus, mockService, mockTelegramCfg)

		if action.Action != cmd.NotifyDown {
			t.Errorf("Expected UP -> DOWN to be NotifyDown, got %v", action.Action)
		}
	})

	// --- SCENARIO 3: Still DOWN (No Action / No transition) ---
	t.Run("Still_DOWN_NoAction", func(t *testing.T) {

		action := sm.ProcessStatus(downStatus, mockService, mockTelegramCfg)

		if action.Action != cmd.NoAction {
			t.Errorf("Expected continuous DOWN check to be NoAction, got %v", action.Action)
		}
	})

	// --- SCENARIO 4: DOWN -> UP Transition (Recovery) ---
	t.Run("DOWN_to_UP_NotifyRecovery_and_Downtime", func(t *testing.T) {
		downtimeDuration := testCheckInterval + 100*time.Millisecond
		time.Sleep(downtimeDuration)

		
		action := sm.ProcessStatus(upStatus, mockService, mockTelegramCfg)

		if action.Action != cmd.NotifyRecovery {
			t.Fatalf("Expected DOWN -> UP to be NotifyRecovery, got %v", action.Action)
		}
		margin := 500 * time.Millisecond
		expectedMinDowntime := downtimeDuration - margin
		expectedMaxDowntime := downtimeDuration + margin

		if action.Downtime < expectedMinDowntime || action.Downtime > expectedMaxDowntime {
			t.Errorf("Downtime mismatch.\nExpected roughly %v, got %v", downtimeDuration, action.Downtime)
		}
	})

	// --- SCENARIO 5: Stable UP (No Action) ---
	t.Run("Stable_UP_NoAction", func(t *testing.T) {
		
		action := sm.ProcessStatus(upStatus, mockService, mockTelegramCfg)

		if action.Action != cmd.NoAction {
			t.Errorf("Expected stable UP check to be NoAction, got %v", action.Action)
		}
	})
}