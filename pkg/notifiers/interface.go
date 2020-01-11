package notifiers

import "context"

// Notifier interface is the interface that
// all notifiers has to implement
type Notifier interface {
	Initialize(context.Context) error
	Configure(map[string]interface{}) error
	Notify(string) error
}
