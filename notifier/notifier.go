package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sentinel/checker"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	// StatusChange indicates a service status change notification
	StatusChange NotificationType = "status_change"
	// ServiceDown indicates a service down notification
	ServiceDown NotificationType = "service_down"
	// ServiceRecovered indicates a service recovery notification
	ServiceRecovered NotificationType = "service_recovered"
)

// Notification represents a notification to be sent
type Notification struct {
	Type        NotificationType       `json:"type"`
	ServiceName string                 `json:"service_name"`
	OldStatus   *checker.ServiceStatus `json:"old_status,omitempty"`
	NewStatus   checker.ServiceStatus  `json:"new_status"`
	Timestamp   time.Time              `json:"timestamp"`
}

// WebhookConfig represents a webhook configuration
type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// DiscordWebhook represents a Discord webhook payload
type DiscordWebhook struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarURL string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed object
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"` // Color in decimal
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"` // ISO8601 timestamp
}

// DiscordEmbedField represents a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// SlackWebhook represents a Slack webhook payload
type SlackWebhook struct {
	Text        string        `json:"text,omitempty"`
	Username    string        `json:"username,omitempty"`
	IconURL     string        `json:"icon_url,omitempty"`
	Attachments []SlackAttach `json:"attachments,omitempty"`
}

// SlackAttach represents a Slack attachment
type SlackAttach struct {
	Fallback   string       `json:"fallback,omitempty"`
	Color      string       `json:"color,omitempty"` // good, warning, danger, or hex
	Pretext    string       `json:"pretext,omitempty"`
	Title      string       `json:"title,omitempty"`
	Text       string       `json:"text,omitempty"`
	Fields     []SlackField `json:"fields,omitempty"`
	Footer     string       `json:"footer,omitempty"`
	FooterIcon string       `json:"footer_icon,omitempty"`
	Timestamp  int64        `json:"ts,omitempty"` // Unix timestamp
}

// SlackField represents a field in a Slack attachment
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

// Notifier interface defines methods for sending notifications
type Notifier interface {
	SendNotification(notification Notification) error
}

// WebhookNotifier implements the Notifier interface for generic webhooks
type WebhookNotifier struct {
	Config WebhookConfig
	Client *http.Client
}

// NewWebhookNotifier creates a new WebhookNotifier
func NewWebhookNotifier(config WebhookConfig) *WebhookNotifier {
	return &WebhookNotifier{
		Config: config,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendNotification sends a notification to the configured webhook
func (n *WebhookNotifier) SendNotification(notification Notification) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshaling notification: %w", err)
	}

	req, err := http.NewRequest("POST", n.Config.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range n.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}

// DiscordNotifier implements the Notifier interface for Discord webhooks
type DiscordNotifier struct {
	Config WebhookConfig
	Client *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier
func NewDiscordNotifier(config WebhookConfig) *DiscordNotifier {
	return &DiscordNotifier{
		Config: config,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendNotification sends a notification to a Discord webhook
func (n *DiscordNotifier) SendNotification(notification Notification) error {
	// Create Discord webhook payload
	var color int
	var title string
	var description string

	switch notification.Type {
	case ServiceDown:
		color = 16711680 // Red
		title = "🔴 Service Down"
		description = fmt.Sprintf("Service **%s** is DOWN", notification.ServiceName)
	case ServiceRecovered:
		color = 65280 // Green
		title = "🟢 Service Recovered"
		description = fmt.Sprintf("Service **%s** is back UP", notification.ServiceName)
	case StatusChange:
		if notification.NewStatus.IsUp {
			color = 65280 // Green
			title = "🟢 Service Status Change"
			description = fmt.Sprintf("Service **%s** changed from DOWN to UP", notification.ServiceName)
		} else {
			color = 16711680 // Red
			title = "🔴 Service Status Change"
			description = fmt.Sprintf("Service **%s** changed from UP to DOWN", notification.ServiceName)
		}
	}

	// Create fields with details
	fields := []DiscordEmbedField{
		{
			Name:   "Service",
			Value:  notification.ServiceName,
			Inline: true,
		},
		{
			Name:   "URL",
			Value:  notification.NewStatus.URL,
			Inline: true,
		},
		{
			Name:   "Status",
			Value:  fmt.Sprintf("%v", notification.NewStatus.IsUp),
			Inline: true,
		},
	}

	// Add response time if available
	if notification.NewStatus.ResponseTime > 0 {
		fields = append(fields, DiscordEmbedField{
			Name:   "Response Time",
			Value:  fmt.Sprintf("%d ms", notification.NewStatus.ResponseTime.Milliseconds()),
			Inline: true,
		})
	}

	// Add status code if available
	if notification.NewStatus.StatusCode > 0 {
		fields = append(fields, DiscordEmbedField{
			Name:   "Status Code",
			Value:  fmt.Sprintf("%d", notification.NewStatus.StatusCode),
			Inline: true,
		})
	}

	// Add error if available
	if notification.NewStatus.Error != nil {
		fields = append(fields, DiscordEmbedField{
			Name:   "Error",
			Value:  notification.NewStatus.Error.Error(),
			Inline: false,
		})
	}

	webhook := DiscordWebhook{
		Username: "SENTINEL Monitoring",
		Embeds: []DiscordEmbed{
			{
				Title:       title,
				Description: description,
				Color:       color,
				Fields:      fields,
				Timestamp:   notification.Timestamp.Format(time.RFC3339),
			},
		},
	}

	payload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("error marshaling Discord webhook: %w", err)
	}

	req, err := http.NewRequest("POST", n.Config.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range n.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Discord webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}

// SlackNotifier implements the Notifier interface for Slack webhooks
type SlackNotifier struct {
	Config WebhookConfig
	Client *http.Client
}

// NewSlackNotifier creates a new SlackNotifier
func NewSlackNotifier(config WebhookConfig) *SlackNotifier {
	return &SlackNotifier{
		Config: config,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendNotification sends a notification to a Slack webhook
func (n *SlackNotifier) SendNotification(notification Notification) error {
	// Create Slack webhook payload
	var color string
	var title string
	var text string

	switch notification.Type {
	case ServiceDown:
		color = "danger" // Red
		title = "🔴 Service Down"
		text = fmt.Sprintf("Service *%s* is DOWN", notification.ServiceName)
	case ServiceRecovered:
		color = "good" // Green
		title = "🟢 Service Recovered"
		text = fmt.Sprintf("Service *%s* is back UP", notification.ServiceName)
	case StatusChange:
		if notification.NewStatus.IsUp {
			color = "good" // Green
			title = "🟢 Service Status Change"
			text = fmt.Sprintf("Service *%s* changed from DOWN to UP", notification.ServiceName)
		} else {
			color = "danger" // Red
			title = "🔴 Service Status Change"
			text = fmt.Sprintf("Service *%s* changed from UP to DOWN", notification.ServiceName)
		}
	}

	// Create fields with details
	fields := []SlackField{
		{
			Title: "Service",
			Value: notification.ServiceName,
			Short: true,
		},
		{
			Title: "URL",
			Value: notification.NewStatus.URL,
			Short: true,
		},
		{
			Title: "Status",
			Value: fmt.Sprintf("%v", notification.NewStatus.IsUp),
			Short: true,
		},
	}

	// Add response time if available
	if notification.NewStatus.ResponseTime > 0 {
		fields = append(fields, SlackField{
			Title: "Response Time",
			Value: fmt.Sprintf("%d ms", notification.NewStatus.ResponseTime.Milliseconds()),
			Short: true,
		})
	}

	// Add status code if available
	if notification.NewStatus.StatusCode > 0 {
		fields = append(fields, SlackField{
			Title: "Status Code",
			Value: fmt.Sprintf("%d", notification.NewStatus.StatusCode),
			Short: true,
		})
	}

	// Add error if available
	if notification.NewStatus.Error != nil {
		fields = append(fields, SlackField{
			Title: "Error",
			Value: notification.NewStatus.Error.Error(),
			Short: false,
		})
	}

	webhook := SlackWebhook{
		Username: "SENTINEL Monitoring",
		Attachments: []SlackAttach{
			{
				Fallback:  title + ": " + text,
				Color:     color,
				Title:     title,
				Text:      text,
				Fields:    fields,
				Footer:    "SENTINEL Monitoring System",
				Timestamp: notification.Timestamp.Unix(),
			},
		},
	}

	payload, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf("error marshaling Slack webhook: %w", err)
	}

	req, err := http.NewRequest("POST", n.Config.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range n.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := n.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Slack webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Slack webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}

// NotifierFactory creates a notifier based on the webhook URL
func NotifierFactory(config WebhookConfig) Notifier {
	// Simple detection based on URL
	if contains(config.URL, "discord.com/api/webhooks") {
		return NewDiscordNotifier(config)
	} else if contains(config.URL, "hooks.slack.com") {
		return NewSlackNotifier(config)
	}

	// Default to generic webhook
	return NewWebhookNotifier(config)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
