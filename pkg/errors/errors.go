package errors

import "fmt"

type Code string

const (
	Invalid       Code = "INVALID_ARGUMENT"
	NotFound      Code = "NOT_FOUND"
	Conflict      Code = "CONFLICT"
	Internal      Code = "INTERNAL"
	Unimplemented Code = "UNIMPLEMENTED"
)

type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return string(e.Code) + ": " + e.Message
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
}
func (e *Error) Unwrap() error {
	if e.Cause == nil {
		return nil
	}
	return fmt.Errorf("%v", e.Cause)
}
func New(c Code, m string) *Error             { return &Error{Code: c, Message: m} }
func Wrap(c Code, m string, err error) *Error { return &Error{Code: c, Message: m, Cause: err} }
