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

// BookApiError is a extension of error containing specific errors within API.
type BookApiError struct {
	errorType OperationError
	msg       string
}

// NewAlreadyCheckedOutError returns a domain already checked out error describing the error.
func NewAlreadyCheckedOutError(id string) *BookApiError {
	return &BookApiError{AlreadyCheckedOut, fmt.Sprintf("Book %s is already checked out", id)}
}

// NewAlreadyCheckedInError returns a domain already checked in error describing the error.
func NewAlreadyCheckedInError(id string) *BookApiError {
	return &BookApiError{AlreadyCheckedIn, fmt.Sprintf("Book %s is already checked in", id)}
}

// NewDatabaseOperationError returns a database connection error describing the error.
func NewDatabaseOperationError(text string) *BookApiError {
	return &BookApiError{DbConnectionError, text}
}

func NewNotFoundError(text string) *BookApiError {
	return &BookApiError{NotFoundError, fmt.Sprintf("Book %s does not exist", text)}
}

// NewMissingEnvVariable returns an missing env variable error describing the error.
func NewMissingEnvVariable(text string) *BookApiError {
	return &BookApiError{MissingEnvVariable, text}
}

// NewValidationError return a domain validation error describing the error
func NewValidationError(text string) *BookApiError {
	return &BookApiError{ValidationError, text}
}

// NewUpdateError returns a domain update error describing the error
func NewUpdateError(text string) *BookApiError {
	return &BookApiError{UpdateError, text}
}

// NewAlreadyExistsError return a domain already exists error
func NewAlreadyExistsError() *BookApiError {
	return &BookApiError{ExistingRecord, "Book already exists"}
}

// NewPersistError returns a domain persistence error
func NewPersistError(text string) *BookApiError {
	return &BookApiError{PersistError, fmt.Sprintf("Error saving domain %s", text)}
}

// Implicit implement of Error interface
func (err *BookApiError) Error() string {
	return err.msg
}
