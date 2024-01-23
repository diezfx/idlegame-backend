package masterdata

import "github.com/diezfx/idlegame-backend/pkg/duration"

type Job struct {
	ID               string            `json:"id"`
	JobType          JobType           `json:"jobType"`
	LevelRequirement int               `json:"levelRequirement"`
	Duration         duration.Duration `json:"duration"`
	Rewards          Reward            `json:"rewards"`
}

func (t JobType) String() string {
	return string(t)
}

type JobContainer struct {
	GatheringJobs  []Job      `json:"gatheringJobs"`
	ProcessingJobs []*Recipes `json:"processingJobs"`
}

// what kind of wood
type JobType string

const (
	WoodcuttingJobType JobType = "woodcutting"
	MiningJobType      JobType = "mining"
	HarvestingJobType  JobType = "harvesting"
	SmeltingJobType    JobType = "smelting"
)

type Reward struct {
	Items      []ItemWithQuantity `json:"items"`
	Experience int                `json:"experience"`
}

func (c *JobContainer) GetGatheringJobDefinition(jobType JobType, jobDefID string) *Job {

	for _, def := range c.GatheringJobs {
		if def.ID == jobDefID && def.JobType == jobType {
			return &def
		}
	}

	return nil
}

func (c *JobContainer) GetSmeltingJobDefinition(jobDefID string) *Recipes {
	for _, def := range c.ProcessingJobs {
		if def.ID == jobDefID && def.JobType == SmeltingJobType {
			return def
		}
	}
	return nil
}

type Recipes struct {
	Job
	Ingredients []ItemWithQuantity `json:"ingredients"`
	OutputItem  ItemWithQuantity   `json:"outputItem"`
}

type ItemWithQuantity struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
