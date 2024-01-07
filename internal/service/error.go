package service

import "errors"

var (
	ErrProjectNotFound        = errors.New("project not found")
	ErrJobTypeNotFound        = errors.New("jobType was not found")
	ErrAlreadyStartedJob      = errors.New("another job was already started")
	ErrLevelRequirementNotMet = errors.New("the level requirement was not met")
)
