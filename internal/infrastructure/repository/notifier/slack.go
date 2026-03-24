package notifier

import (
	"bytes"
	"fmt"
	"net/http"
)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Notify(message string) error {
	// Create JSON payload
	payload := fmt.Sprintf(`{"text":"%s"}`, message)

	// Send POST request to Slack webhook URL
	_, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer([]byte(payload)))
	return err
}
