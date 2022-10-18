package errors

import "fmt"

type Unknown struct {
	Tag   string
	Cause error
}

func (e *Unknown) Error() string {
	return fmt.Sprintf("[%s] Unknown error. Caused by: %s", e.Tag, e.Cause.Error())
}

type InvalidField struct {
	Domain string
	Field  string
	Value  interface{}
}

func (e *InvalidField) Error() string {
	return fmt.Sprintf("[%s] Invalid field %s with value %v", e.Domain, e.Field, e.Value)
}

type MultipleInvalidFields struct {
	Errors []error
}

func (e *MultipleInvalidFields) Error() string {
	return fmt.Sprintf("Multiple errors: %v", e.Errors)
}
