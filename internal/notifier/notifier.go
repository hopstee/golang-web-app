package notifier

import "mobile-backend-boilerplate/internal/repository"

type Notifier interface {
	SendMessage(request repository.Request) error
	SendMessageWithRetry(request repository.Request) error
}
