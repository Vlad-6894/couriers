package pkg_errors

import "errors"

var (
	ErrConflict        = errors.New("Conflict Error: ")
	ErrInvalidArgument = errors.New("Invalid Argument Error: ")
	ErrNotFound        = errors.New("Not Found Error: ")
	ErrInvalidPassword = errors.New("Invalid Password Error: ")
)
