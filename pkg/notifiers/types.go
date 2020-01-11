package notifiers

import "fmt"

// FromString returns a `Check` parsed by a string
func FromString(checkType string) (Notifier, error) {
	switch checkType {
	case "slack":
		return &SlackNotifier{}, nil
	}
	return nil, fmt.Errorf("no such type: %s", checkType)
}
