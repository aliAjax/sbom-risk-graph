package attestation

import "fmt"

type Receipt struct {
	State string
	IDs   []string
}

type ImportError struct{ Cause error }

func (e *ImportError) Error() string     { return fmt.Sprintf("attestation import failed: %v", e.Cause) }
func (e *ImportError) CauseError() error { return e.Cause }

func wrapImportError(cause error) error {
	return fmt.Errorf("receive attestation: %v", &ImportError{Cause: cause})
}
