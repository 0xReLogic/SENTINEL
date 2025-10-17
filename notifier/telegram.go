package notifier

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)


var markdownReplacer = strings.NewReplacer(
	"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
	"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">",
	"\\>", "#", "\\#", "+", "\\+", "-", "\\-", "=",
	"\\=", "|", "\\|", "{", "\\{", "}", "\\}", ".",
	"\\.", "!", "\\!",
)

func escapeMarkdownV2(s string) string {
	return markdownReplacer.Replace(s)
}


// SendTelegramNotification sends a message to Telegram using the production API URL.
func SendTelegramNotification(token, chatID, message string) error {
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	return sendTelegramRequest(token, chatID, message, apiUrl)
}

func sendTelegramRequest(token, chatID, message, apiUrl string) error {
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


// (FormatDownMessage and FormatRecoveryMessage functions remain the same)
func FormatDownMessage(name, url, errorMsg string, checkTime time.Time) string {
	return fmt.Sprintf("🔴 *Service DOWN*\n*Name:* %s\n*URL:* %s\n*Error:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(errorMsg),
		escapeMarkdownV2(checkTime.Format("2006-01-02 15:04:05")),
	)
}

func FormatRecoveryMessage(name, url string, downtime time.Duration, recoveryTime time.Time) string {
	return fmt.Sprintf("🟢 *Service RECOVERED*\n*Name:* %s\n*URL:* %s\n*Downtime:* %s\n*Time:* %s",
		escapeMarkdownV2(name),
		escapeMarkdownV2(url),
		escapeMarkdownV2(downtime.String()),
		escapeMarkdownV2(recoveryTime.Format("2006-01-02 15:04:05")),
	)
}