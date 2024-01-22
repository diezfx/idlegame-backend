package jobs

import (
	"time"

	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

type Job struct {
	ID        int
	JobDefID  string
	UserID    int
	StartedAt time.Time
	UpdatedAt time.Time
	Monsters  []int
	JobType   string
}

func FromWoodcuttingJob(j *storage.Job) *WoodCuttingJob {
	return &WoodCuttingJob{
		Job:      fromJob(*j),
		TreeType: item.TreeType(j.JobDefID),
	}
}

func FromMiningJob(j *storage.Job) *MiningJob {
	return &MiningJob{
		Job:     fromJob(*j),
		OreType: item.OreType(j.JobDefID),
	}
}

func FromHarvestingJob(j *storage.Job) *HarvestingJob {
	return &HarvestingJob{
		Job:      fromJob(*j),
		CropType: item.CropType(j.JobDefID),
	}
}

type WoodCuttingJob struct {
	Job
	TreeType item.TreeType `json:"treeType"`
}

type MiningJob struct {
	Job
	OreType item.OreType `json:"oreType"`
}

type HarvestingJob struct {
	Job
	CropType item.CropType `json:"cropType"`
}
