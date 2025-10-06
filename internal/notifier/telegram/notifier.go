package telegram

import (
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/notifier"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/pkg/helper/markdown"
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

func (n *TelegramNotifier) SendMessage(request repository.Request) error {
	n.Logger.Info("send telegram notification message attempt")
	apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.token)

	msg := fmt.Sprintf("*Имя:* %s\n", markdown.EscapeMarkdownV2(request.Name))
	msg += fmt.Sprintf("*Сообщение:* %s\n", markdown.EscapeMarkdownV2(request.Message))
	msg += fmt.Sprintf("*Способ связи:* %s\n", markdown.EscapeMarkdownV2(request.ContactType))
	msg += fmt.Sprintf("*Телефон:* %s\n", markdown.EscapeMarkdownV2(request.Phone))
	msg += fmt.Sprintf("*Почта:* %s\n", markdown.EscapeMarkdownV2(request.Email))

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

func (n *TelegramNotifier) SendMessageWithRetry(request repository.Request) error {
	return nil
}
