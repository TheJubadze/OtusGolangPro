package errors

import "fmt"

type Error interface {
	error
	Unwrap() error
}

type ProgramError struct {
	Err error
}

func (e *ProgramError) Error() string {
	return fmt.Sprintf("Program error: %v", e.Err)
}

func (e *ProgramError) Unwrap() error {
	return e.Err
}

type TagError struct {
	Key   string
	Value string
	Err   error
}

func (e *TagError) Error() string {
	return fmt.Sprintf("Tag error: invalid %s value: %s", e.Key, e.Value)
}

func (e *TagError) Unwrap() error {
	return e.Err
}
