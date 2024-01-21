package api

import (
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service/jobs"
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

type StartWoodCuttingJobRequest struct {
	UserID   int           `json:"userId"`
	Monster  int           `json:"monster"`
	TreeType jobs.TreeType `json:"treeType"`
}

type StartMiningJobRequest struct {
	UserID  int          `json:"userId"`
	Monster int          `json:"monster"`
	OreType jobs.OreType `json:"oreType"`
}

type StartHarvestingJobRequest struct {
	UserID   int           `json:"userId"`
	Monster  int           `json:"monster"`
	CropType jobs.CropType `json:"cropType"`
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
		Type:       m.Type().String(),
		Level:      m.Level(),
		Experience: m.Experience,
	}
}
