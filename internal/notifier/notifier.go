package notifier

type Notifier interface {
	SendMessage(msg string) error
}
