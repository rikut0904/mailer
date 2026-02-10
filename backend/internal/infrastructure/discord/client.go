package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	webhookURL string
}

func NewClient(webhookURL string) *Client {
	return &Client{webhookURL: webhookURL}
}

type WebhookPayload struct {
	Content string `json:"content"`
}

func (c *Client) SendNotification(message string) error {
	if c.webhookURL == "" {
		return nil
	}

	payload := WebhookPayload{Content: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
