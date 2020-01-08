package checks

import "fmt"

// FromString returns a `Check` parsed by a string
func FromString(checkType string) (Check, error) {
	switch checkType {
	case "http":
		return &HTTPCheck{}, nil
	}
	return nil, fmt.Errorf("no such type: %s", checkType)
}
