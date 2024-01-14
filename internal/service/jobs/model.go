package jobs

import (
	"time"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
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
	Job
	TreeType TreeType `json:"treeType"`
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
	UserID    int
	StartedAt time.Time
	UpdatedAt time.Time
	Monsters  []int
	JobType   string
}

func FromWoodcuttingJob(j *storage.WoodCuttingJob) *WoodCuttingJob {
	return &WoodCuttingJob{
		Job: Job{
			ID:        j.ID,
			UserID:    j.UserID,
			StartedAt: j.StartedAt,
			Monsters:  j.Monsters,
			JobType:   WoodCuttingJobType.String(),
			UpdatedAt: j.UpdatedAt,
		},
		TreeType: TreeType(j.TreeType),
	}
}

type Reward struct {
	Items []inventory.Item
	Exp   int
}
