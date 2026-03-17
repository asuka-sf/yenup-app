package notifier

type Notifier interface {
	// Notify the user via　Slack
	Notify(message string) error
}
