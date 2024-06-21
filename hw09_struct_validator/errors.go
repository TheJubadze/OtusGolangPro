package hw09structvalidator

import "fmt"

type Error interface {
	error
	Unwrap() error
	IsValidation() bool
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

func (e *ProgramError) IsValidation() bool {
	return false
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

func (e *TagError) IsValidation() bool {
	return false
}

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: invalid field '%s': %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func (e *ValidationError) IsValidation() bool {
	return true
}
