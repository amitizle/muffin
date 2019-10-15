package checks

// Check interface is the interface that
// all checks has to implement
type Check interface {
	Run() (string, error)
}
