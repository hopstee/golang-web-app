package infrastructure

import "mobile-backend-boilerplate/internal/notifier/telegram"

func (d *Dependencies) InitNotifiers() {
	d.TelegramNotifier = telegram.NewTelegramNotifier(d.Config.Telegram.Token, d.Config.Telegram.ChatID, d.Logger)
}
