package apperror

import "errors"

type AuthenticationError struct{ msg string }
type ForbiddenError struct{ msg string }
type NotFoundError struct{ msg string }
type DuplicateError struct{ msg string }
type InsertError struct{ msg string }
type UpdateError struct{ msg string }
type DomainValidationError struct{ msg string }

func (e *AuthenticationError) Error() string   { return e.msg }
func (e *ForbiddenError) Error() string        { return e.msg }
func (e *NotFoundError) Error() string         { return e.msg }
func (e *DuplicateError) Error() string        { return e.msg }
func (e *InsertError) Error() string           { return e.msg }
func (e *UpdateError) Error() string           { return e.msg }
func (e *DomainValidationError) Error() string { return e.msg }

func NewAuthenticationError(msg string) *AuthenticationError { return &AuthenticationError{msg} }
func NewForbiddenError(msg string) *ForbiddenError           { return &ForbiddenError{msg} }
func NewNotFoundError(msg string) *NotFoundError             { return &NotFoundError{msg} }
func NewDuplicateError(msg string) *DuplicateError           { return &DuplicateError{msg} }
func NewInsertError(msg string) *InsertError                 { return &InsertError{msg} }
func NewUpdateError(msg string) *UpdateError                 { return &UpdateError{msg} }
func NewDomainValidationError(msg string) *DomainValidationError {
	return &DomainValidationError{msg}
}

func IsAuthenticationError(err error) bool {
	var e *AuthenticationError
	return errors.As(err, &e)
}

func IsForbiddenError(err error) bool {
	var e *ForbiddenError
	return errors.As(err, &e)
}

func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

func IsDuplicateError(err error) bool {
	var e *DuplicateError
	return errors.As(err, &e)
}

func IsInsertError(err error) bool {
	var e *InsertError
	return errors.As(err, &e)
}

func IsUpdateError(err error) bool {
	var e *UpdateError
	return errors.As(err, &e)
}
