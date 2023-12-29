package api

import (
	"fmt"
)

type InvalidArgumentError struct {
	Argument string
}

func (a *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument %s", a.Argument)
}

func NewInvalidArgumentError(arg string) *InvalidArgumentError {
	return &InvalidArgumentError{Argument: arg}
}

type ErrorResponse struct {
	ErrorCode int
	Reason    string
}
