package telegram

import (
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/notifier"
	"net/http"
	"net/url"
)

type TelegramNotifier struct {
	token  string
	chatID string
	Logger *slog.Logger
}

func NewTelegramNotifier(token string, chatID string, logger *slog.Logger) notifier.Notifier {
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		Logger: logger,
	}
}

func (n *TelegramNotifier) SendMessage(msg string) error {
	n.Logger.Info("send telegram notification message attempt")
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.token)

	resp, err := http.PostForm(apiUrl, url.Values{
		"chat_id":    {n.chatID},
		"text":       {msg},
		"parse_mode": {"MarkdownV2"},
	})
	if err != nil {
		n.Logger.Error("send telegram notification message: failed to send message", slog.Any("err", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		n.Logger.Error("send telegram notification message: telegram api error", slog.Any("err", err))
		return fmt.Errorf("telegram api error: %v", resp.Status)
	}

	return nil
}
