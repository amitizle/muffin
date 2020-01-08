package checks

import "context"

// Check interface is the interface that
// all checks has to implement
type Check interface {
	Initialize(context.Context) error
	Configure(map[string]interface{}) error
	Run() ([]byte, error)
}
