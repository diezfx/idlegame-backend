package jobs

import (
	"time"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

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

type JobContainer struct {
	woodcuttingDefs []JobDefinition
	miningDefs      []JobDefinition
	harvestingDefs  []JobDefinition

	smeltingDefs []*Recipes
}

// what kind of wood
type JobType string

const (
	WoodCuttingJobType JobType = "WoodCutting"
	MiningJobType      JobType = "Mining"
	HarvestingJobType  JobType = "Harvesting"
	SmeltingJobType    JobType = "Smelting"
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

type Reward struct {
	Items []inventory.Item
	Exp   int
}

type JobDefinition struct {
	JobDefID         string
	JobType          JobType
	LevelRequirement int
	Duration         time.Duration
	Rewards          Reward
}

func (c *JobContainer) GetGatheringJobDefinition(jobType JobType, jobDefID string) *JobDefinition {
	if jobType == WoodCuttingJobType {
		for _, def := range c.woodcuttingDefs {
			if def.JobDefID == jobDefID {
				return &def
			}
		}
	}

	if jobType == MiningJobType {
		for _, def := range c.miningDefs {
			if def.JobDefID == jobDefID {
				return &def
			}
		}
	}

	if jobType == HarvestingJobType {
		for _, def := range c.harvestingDefs {
			if def.JobDefID == jobDefID {
				return &def
			}
		}
	}
	return nil
}

func (c *JobContainer) GetSmeltingJobDefinition(jobDefID string) *Recipes {
	for _, def := range c.smeltingDefs {
		if def.JobDefID == jobDefID {
			return def
		}
	}
	return nil
}
