package jobs

import (
	"time"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

// what kind of wood
// different tree types have
// level requirements, durations, exp gains
type TreeType string

const (
	SpruceType TreeType = "Spruce"
	BirchType  TreeType = "Birch"
	PineType   TreeType = "Pine"
)

func (t TreeType) String() string {
	return string(t)
}

type WoodCuttingJob struct {
	ID        int       `json:"id"`
	Monster   int       `json:"monster"`
	TreeType  TreeType  `json:"treeType"`
	StartedAt time.Time `json:"startedAt"`
}

type WoodCuttingJobContainer struct {
	defs []WoodCuttingJobDefinition
}

// what kind of wood
type JobType string

const (
	WoodCuttingJobType JobType = "WoodCutting"
)

type Job struct {
	ID        int
	StartedAt time.Time
	Monsters  []int
	JobType   string
}

func FromWoodcuttingJob(j *storage.WoodCuttingJob) *WoodCuttingJob {
	return &WoodCuttingJob{
		ID:        j.ID,
		Monster:   j.MonsterID,
		TreeType:  TreeType(j.TreeType),
		StartedAt: j.StartedAt,
	}
}
