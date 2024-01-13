package api

import (
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service/jobs"
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

type StartWoodCuttingJobRequest struct {
	UserID   int           `json:"userId"`
	Monster  int           `json:"monster"`
	TreeType jobs.TreeType `json:"treeType"`
}
