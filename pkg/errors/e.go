package errors

import "fmt"

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(temp string, args ...any) error {
	return &errorString{
		s: fmt.Sprintf(temp, args...),
	}
}
