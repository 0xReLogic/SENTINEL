// Package notifier handles the sending of notifications to external services like Telegram.
// It is responsible for formatting messages and interacting with provider APIs.
package notifier

import (
	"fmt"
	"net/http"
    "io"
	"net/url"
	"strings"
	"time"
)

// markdownReplacer is a pre-built, reusable strings.Replacer for escaping Telegram's MarkdownV2 characters.
// Defining it once at the package level is much more efficient than creating it on every function call.
var markdownReplacer = strings.NewReplacer(
	"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
	"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">",
	"\\>", "#", "\\#", "+", "\\+", "-", "\\-", "=",
	"\\=", "|", "\\|", "{", "\\{", "}", "\\}", ".",
	"\\.", "!", "\\!",
)

// escapeMarkdownV2 escapes characters that have special meaning in Telegram's MarkdownV2 format.
// It uses the pre-built, package-level markdownReplacer for efficiency.
func escapeMarkdownV2(s string) string {
	return markdownReplacer.Replace(s)
}

// SendTelegramNotification sends a formatted message to a specified Telegram chat.
// It uses a custom HTTP client to ensure all requests have a timeout.
func SendTelegramNotification(token, chatID, message string) error {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	params := url.Values{}
	params.Set("chat_id", chatID)
	params.Set("text", message)
	params.Set("parse_mode", "MarkdownV2")

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use a client with a timeout to prevent the application from hanging.
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API returned status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// FormatDownMessage creates the standardized notification text for when a service goes down.
// It uses escapeMarkdownV2 on all dynamic parts of the message to ensure it is rendered correctly by Telegram.
func FormatDownMessage(name, url, errorMsg string, checkTime time.Time) string {
	return fmt.Sprintf("🔴 *Service DOWN*\n*Name:* %s\n*URL:* %s\n*Error:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(errorMsg),
		escapeMarkdownV2(checkTime.Format("2006-01-02 15:04:05")),
	)
}

// FormatRecoveryMessage creates the standardized notification text for when a service recovers.
// All dynamic parts of the message are escaped using escapeMarkdownV2 to prevent formatting issues.
func FormatRecoveryMessage(name, url string, downtime time.Duration, recoveryTime time.Time) string {
	return fmt.Sprintf("🟢 *Service RECOVERED*\n*Name:* %s\n*URL:* %s\n*Downtime:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(downtime.String()),
		escapeMarkdownV2(recoveryTime.Format("2006-01-02 15:04:05")),
	)
}