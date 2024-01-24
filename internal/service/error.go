package service

import "errors"

var (
	ErrJobNotFound            = errors.New("job not found")
	ErrJobTypeNotFound        = errors.New("jobType was not found")
	ErrMonsterNotFound        = errors.New("monster not found")
	ErrAlreadyStartedJob      = errors.New("another job was already started")
	ErrLevelRequirementNotMet = errors.New("the level requirement was not met")
	ErrNotEnoughInventory     = errors.New("not enough inventory")
	ErrNotEnoughItems         = errors.New("not enough items")
)
