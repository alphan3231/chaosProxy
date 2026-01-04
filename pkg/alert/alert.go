package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payload struct {
	Text string `json:"text"`
}

func Send(webhookURL, message string) error {
	if webhookURL == "" {
		return nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	payload := Payload{
		Text: fmt.Sprintf("👻 **Sentinel Alert**\n%s", message),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("alert failed with status: %d", resp.StatusCode)
	}

	return nil
}
