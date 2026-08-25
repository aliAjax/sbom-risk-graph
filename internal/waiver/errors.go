package waiver

import "fmt"

type ReviewError struct {
	WaiverID string
	Cause    error
}

func (e *ReviewError) Error() string {
	return fmt.Sprintf("waiver %s review failed: %v", e.WaiverID, e.Cause)
}

func (e *ReviewError) Unwrap() error { return e.Cause }

func wrapReviewError(id string, cause error) error {
	return fmt.Errorf("review decision: %w", &ReviewError{WaiverID: id, Cause: cause})
}
