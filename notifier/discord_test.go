package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFormatDownEmbed(t *testing.T) {
	name := "Test Service"
	url := "https://api.test-service.com"
	errorMsg := "connection timeout"
	checkTime := time.Date(2025, 10, 11, 22, 30, 0, 0, time.UTC)

	embed := FormatDownEmbed(name, url, errorMsg, checkTime)

	if embed.Title != "🔴 Service DOWN" {
		t.Errorf("Expected title '🔴 Service DOWN', got '%s'", embed.Title)
	}

	if embed.Color != ColorRed {
		t.Errorf("Expected color %d, got %d", ColorRed, embed.Color)
	}

	if len(embed.Fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(embed.Fields))
	}

	if embed.Fields[0].Name != "Service" || embed.Fields[0].Value != name {
		t.Errorf("Expected Service field with value '%s', got '%s'", name, embed.Fields[0].Value)
	}

	if embed.Fields[1].Name != "URL" || embed.Fields[1].Value != url {
		t.Errorf("Expected URL field with value '%s', got '%s'", url, embed.Fields[1].Value)
	}

	if embed.Fields[2].Name != "Error" || embed.Fields[2].Value != errorMsg {
		t.Errorf("Expected Error field with value '%s', got '%s'", errorMsg, embed.Fields[2].Value)
	}

	expectedTimestamp := "2025-10-11T22:30:00Z"
	if embed.Timestamp != expectedTimestamp {
		t.Errorf("Expected timestamp '%s', got '%s'", expectedTimestamp, embed.Timestamp)
	}
}

func TestFormatRecoveryEmbed(t *testing.T) {
	name := "Test Service"
	url := "https://api.test-service.com"
	downtime := 5*time.Minute + 30*time.Second
	recoveryTime := time.Date(2025, 10, 11, 22, 35, 30, 0, time.UTC)

	embed := FormatRecoveryEmbed(name, url, downtime, recoveryTime)

	if embed.Title != "🟢 Service RECOVERED" {
		t.Errorf("Expected title '🟢 Service RECOVERED', got '%s'", embed.Title)
	}

	if embed.Color != ColorGreen {
		t.Errorf("Expected color %d, got %d", ColorGreen, embed.Color)
	}

	if len(embed.Fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(embed.Fields))
	}

	if embed.Fields[0].Name != "Service" || embed.Fields[0].Value != name {
		t.Errorf("Expected Service field with value '%s', got '%s'", name, embed.Fields[0].Value)
	}

	if embed.Fields[1].Name != "URL" || embed.Fields[1].Value != url {
		t.Errorf("Expected URL field with value '%s', got '%s'", url, embed.Fields[1].Value)
	}

	if embed.Fields[2].Name != "Downtime" || embed.Fields[2].Value != "5m30s" {
		t.Errorf("Expected Downtime field with value '5m30s', got '%s'", embed.Fields[2].Value)
	}

	expectedTimestamp := "2025-10-11T22:35:30Z"
	if embed.Timestamp != expectedTimestamp {
		t.Errorf("Expected timestamp '%s', got '%s'", expectedTimestamp, embed.Timestamp)
	}
}

func TestSendDiscordNotification_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		var payload DiscordWebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Failed to decode payload: %v", err)
		}

		if payload.Username != "SENTINEL Monitor" {
			t.Errorf("Expected username 'SENTINEL Monitor', got '%s'", payload.Username)
		}

		if len(payload.Embeds) != 1 {
			t.Fatalf("Expected 1 embed, got %d", len(payload.Embeds))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	embed := FormatDownEmbed("Test", "https://test.com", "error", time.Now())
	err := SendDiscordNotification(server.URL, "", embed)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestSendDiscordNotification_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid webhook"}`))
	}))
	defer server.Close()

	embed := FormatDownEmbed("Test", "https://test.com", "error", time.Now())
	err := SendDiscordNotification(server.URL, "", embed)
	if err == nil {
		t.Fatal("Expected an error for failed API call, but got nil")
	}

	if !strings.Contains(err.Error(), "Discord API returned status code 400") {
		t.Errorf("Expected error to contain 'Discord API returned status code 400', got: %v", err)
	}
}

func TestSendDiscordNotification_NetworkTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	embed := FormatDownEmbed("Test", "https://test.com", "error", time.Now())
	err := SendDiscordNotification(server.URL, "", embed)
	if err == nil {
		t.Fatal("Expected a timeout error, but got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected error to contain 'context deadline exceeded', got: %v", err)
	}
}
