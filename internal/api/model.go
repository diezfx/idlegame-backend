package api

import (
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service/monster"
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

type StartGatheringJob struct {
	UserID   int    `json:"userId"`
	Monster  int    `json:"monster"`
	JobDefID string `json:"jobDefId"`
}

type StartProcessingJob struct {
	UserID   int    `json:"userId"`
	Monster  int    `json:"monster"`
	JobDefID string `json:"jobDefId"`
}

type Monster struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}

func toMonster(m *monster.Monster) Monster {
	return Monster{
		ID:         m.ID,
		Name:       m.Name(),
		Type:       m.Element().String(),
		Level:      m.Level(),
		Experience: m.Experience,
	}
}
