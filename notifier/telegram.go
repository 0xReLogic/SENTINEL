// Package notifier handles the sending of notifications to external services like Telegram.
// It is responsible for formatting messages and interacting with provider APIs.
package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// escapeMarkdownV2 escapes characters that have special meaning in Telegram's MarkdownV2 format.
// Telegram requires specific characters like '.', '-', '(', ')', etc., to be escaped with a
// preceding backslash to be displayed as literal characters. This function ensures that any user-provided
// string (like a URL or an error message) will not break the markdown parsing on Telegram's side.
func escapeMarkdownV2(s string) string {
	var replacer = strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
		"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">",
		"\\>", "#", "\\#", "+", "\\+", "-", "\\-", "=",
		"\\=", "|", "\\|", "{", "\\{", "}", "\\}", ".",
		"\\.", "!", "\\!",
	)
	return replacer.Replace(s)
}

// SendTelegramNotification sends a formatted message to a specified Telegram chat.
// It constructs the API request, sets the required parameters (chat_id, text, parse_mode),
// and executes an HTTP POST request. It handles potential network errors and checks for a
// successful (200 OK) status code from the Telegram API, returning an error if the
// notification fails to send.
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status code %d", resp.StatusCode)
	}

	return nil
}

// FormatDownMessage creates the standardized notification text for when a service goes down.
// It takes the service details and the specific error as input and formats them into a
// user-friendly, multi-line string. It uses escapeMarkdownV2 on all dynamic parts of the
// message to ensure it is rendered correctly by Telegram.
func FormatDownMessage(name, url, errorMsg string, checkTime time.Time) string {
	return fmt.Sprintf("🔴 *Service DOWN*\n*Name:* %s\n*URL:* %s\n*Error:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(errorMsg),
		escapeMarkdownV2(checkTime.Format("2006-01-02 15:04:05")),
	)
}

// FormatRecoveryMessage creates the standardized notification text for when a service recovers.
// It includes the service name, URL, total downtime, and the time of recovery.
// All dynamic parts of the message are escaped using escapeMarkdownV2 to prevent formatting
// issues in the Telegram client.
func FormatRecoveryMessage(name, url string, downtime time.Duration, recoveryTime time.Time) string {
	return fmt.Sprintf("🟢 *Service RECOVERED*\n*Name:* %s\n*URL:* %s\n*Downtime:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(downtime.String()),
		escapeMarkdownV2(recoveryTime.Format("2006-01-02 15:04:05")),
	)
}
