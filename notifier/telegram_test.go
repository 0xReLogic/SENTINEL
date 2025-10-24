package notifier

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)


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

func TestEscapeMarkdownV2(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {
            name:  "No special chars",
            input: "Simple text",
            want:  "Simple text",
        },
        {
            name:  "All special chars",
            input: "._*[]()~`>#+-=|{}!.!",
            want:  "\\.\\_\\*\\[\\]\\(\\)\\~\\`\\>\\#\\+\\-\\=\\|\\{\\}\\!\\.\\!",
        },
        {
            name:  "Mixed text and chars",
            input: "Service *down* at 12:30!",
            want:  "Service \\*down\\* at 12:30\\!",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := escapeMarkdownV2(tt.input)
            if got != tt.want {
                t.Errorf("escapeMarkdownV2() got = %q, want %q", got, tt.want)
            }
        })
    }
}


func TestSendTelegramNotification_Success(t *testing.T) {
	// Create a mock server that simulates a successful Telegram API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}
		if r.FormValue("chat_id") != "test-chat" {
			t.Errorf("Expected chat_id 'test-chat', got '%s'", r.FormValue("chat_id"))
		}
		if r.FormValue("text") != "Hello World" {
			t.Errorf("Expected text 'Hello World', got '%s'", r.FormValue("text"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	// Call the function with our mock server's URL
	err := sendTelegramRequest("test-token", "test-chat", "Hello World", server.URL)
	if err != nil {
		t.Errorf("Expected no error for a successful send, but got: %v", err)
	}
}

// TestSendTelegramNotification_APIError tests the handling of an error from the Telegram API.
func TestSendTelegramNotification_APIError(t *testing.T) {
	// Create a mock server that simulates a Telegram API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"ok":false, "description":"Bad Request: chat not found"}`))
	}))
	defer server.Close()

	err := sendTelegramRequest("test-token", "invalid-chat", "message", server.URL)
	if err == nil {
		t.Fatal("Expected an error for a failed API call, but got nil")
	}

	expectedErrorMsg := "telegram API returned status code 400: {\"ok\":false, \"description\":\"Bad Request: chat not found\"}"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

// TestSendTelegramNotification_NetworkTimeout tests the handling of a client-side timeout.
func TestSendTelegramNotification_NetworkTimeout(t *testing.T) {
	// Create a server that waits longer than the client's timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second) // Client timeout is 10 seconds
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := sendTelegramRequest("test-token", "test-chat", "message", server.URL)
	if err == nil {
		t.Fatal("Expected a timeout error, but got nil")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected error to contain 'context deadline exceeded', got: %v", err)
	}
}