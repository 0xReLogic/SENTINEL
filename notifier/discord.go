package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Discord embed colors
const (
	ColorRed    = 15158332 // #E74C3C - DOWN status
	ColorGreen  = 3066993  // #2ECC71 - RECOVERY status
	ColorOrange = 15105570 // #E67E22 - Degraded (future use)
)

// DiscordEmbed represents a Discord embed object
type DiscordEmbed struct {
	Title     string         `json:"title"`
	Color     int            `json:"color"`
	Fields    []DiscordField `json:"fields"`
	Timestamp string         `json:"timestamp"`
}

// DiscordField represents a field in a Discord embed
type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordWebhookPayload represents the payload sent to Discord webhook
type DiscordWebhookPayload struct {
	Username string         `json:"username"`
	Embeds   []DiscordEmbed `json:"embeds"`
}

// SendDiscordNotification sends a message to Discord using webhook URL
func SendDiscordNotification(webhookURL, message string, embed DiscordEmbed) error {
	payload := DiscordWebhookPayload{
		Username: "SENTINEL Monitor",
		Embeds:   []DiscordEmbed{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Discord request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Discord API returned status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// FormatDownEmbed creates a Discord embed for service DOWN notification
func FormatDownEmbed(name, url, errorMsg string, checkTime time.Time) DiscordEmbed {
	return DiscordEmbed{
		Title: "🔴 Service DOWN",
		Color: ColorRed,
		Fields: []DiscordField{
			{Name: "Service", Value: name, Inline: true},
			{Name: "URL", Value: url, Inline: true},
			{Name: "Error", Value: errorMsg, Inline: false},
		},
		Timestamp: checkTime.Format(time.RFC3339),
	}
}

// FormatRecoveryEmbed creates a Discord embed for service RECOVERY notification
func FormatRecoveryEmbed(name, url string, downtime time.Duration, recoveryTime time.Time) DiscordEmbed {
	return DiscordEmbed{
		Title: "🟢 Service RECOVERED",
		Color: ColorGreen,
		Fields: []DiscordField{
			{Name: "Service", Value: name, Inline: true},
			{Name: "URL", Value: url, Inline: true},
			{Name: "Downtime", Value: downtime.String(), Inline: false},
		},
		Timestamp: recoveryTime.Format(time.RFC3339),
	}
}
