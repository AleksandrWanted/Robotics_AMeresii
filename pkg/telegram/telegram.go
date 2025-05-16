package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Message struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func SendMessage(text string) error {
	tgEndpoint := os.Getenv("TELEGRAM_ENDPOINT")
	tgBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	tgChatID := os.Getenv("TELEGRAM_TARGET_CHAT_ID")

	message := Message{
		ChatID: tgChatID,
		Text:   text,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", tgEndpoint, tgBotToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", body)
	}

	return nil
}
