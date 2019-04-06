package domain

import "fmt"

type OperationError int8

const (
	AlreadyCheckedOut OperationError = iota
	AlreadyCheckedIn
	UpdateError
	ValidationError
	ExistingRecord
	MissingEnvVariable
	DbConnectionError
	NotFoundError
	PersistError
)

func (oe OperationError) Name() string {
	names := [...]string{
		"AlreadyCheckedOut",
		"AlreadyCheckedIn",
		"UpdateError",
		"ValidationError",
		"ExistingRecord",
		"MissingEnvVariable",
		"DbConnectionError",
		"NotFoundError",
		"PersistError",
	}

	// prevent panicking in case of
	// `status` is out of range
	if oe < AlreadyCheckedIn || oe > PersistError {
		return "Unknown"
	}

	return names[oe]
}

// BookAPIError is a extension of error containing specific errors within API.
type BookAPIError struct {
	errorType OperationError
	msg       string
}

// NewAlreadyCheckedOutError returns a domain already checked out error describing the error.
func NewAlreadyCheckedOutError(id string) *BookAPIError {
	return &BookAPIError{AlreadyCheckedOut, fmt.Sprintf("Book %s is already checked out", id)}
}

// NewAlreadyCheckedInError returns a domain already checked in error describing the error.
func NewAlreadyCheckedInError(id string) *BookAPIError {
	return &BookAPIError{AlreadyCheckedIn, fmt.Sprintf("Book %s is already checked in", id)}
}

// NewDatabaseOperationError returns a database connection error describing the error.
func NewDatabaseOperationError(text string) *BookAPIError {
	return &BookAPIError{DbConnectionError, text}
}

// NewNotFoundError  returns a not found error describing the error.
func NewNotFoundError(text string) *BookAPIError {
	return &BookAPIError{NotFoundError, fmt.Sprintf("Book %s does not exist", text)}
}

// NewMissingEnvVariable returns an missing env variable error describing the error.
func NewMissingEnvVariable(text string) *BookAPIError {
	return &BookAPIError{MissingEnvVariable, text}
}

// NewValidationError return a domain validation error describing the error
func NewValidationError(text string) *BookAPIError {
	return &BookAPIError{ValidationError, text}
}

// NewUpdateError returns a domain update error describing the error
func NewUpdateError(text string) *BookAPIError {
	return &BookAPIError{UpdateError, text}
}

// NewAlreadyExistsError return a domain already exists error
func NewAlreadyExistsError() *BookAPIError {
	return &BookAPIError{ExistingRecord, "Book already exists"}
}

// NewPersistError returns a domain persistence error
func NewPersistError(text string) *BookAPIError {
	return &BookAPIError{PersistError, fmt.Sprintf("Error saving domain %s", text)}
}

// Implicit implement of Error interface
func (err *BookAPIError) Error() string {
	return err.msg
}
