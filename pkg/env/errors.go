package env

import (
	"errors"
)

var (
	// errInvalidValue returned when the value passed to Unmarshal is nil or not a
	// pointer to a struct.
	errInvalidValue = errors.New("value must be a non-nil pointer to a struct")

	// errUnsupportedType returned when a field with tag "env" is unsupported.
	errUnsupportedType = errors.New("field is an unsupported type")

	// errUnexportedField returned when a field with tag "env" is not exported.
	errUnexportedField = errors.New("field must be exported")

	// errInvalidEnviron returned when environ has an incorrect format.
	errInvalidEnviron = errors.New("items in environ must have format key=value")
)
