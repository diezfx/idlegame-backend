package masterdata

import "github.com/diezfx/idlegame-backend/pkg/duration"

type JobType string

type Job struct {
	ID               string            `json:"id"`
	JobType          JobType           `json:"jobType"`
	LevelRequirement int               `json:"levelRequirement"`
	Duration         duration.Duration `json:"duration"`
	Rewards          Reward            `json:"rewards"`
	Affinity         Affinity          `json:"affinity"`
}

// specific element types may be better at some job; this adds a multiplicator. Can be used to reduce time or increase rewards
type Affinity struct {
	Elements map[MonsterElement]float64 `json:"elements"`
}

func (t JobType) String() string {
	return string(t)
}

type JobContainer struct {
	GatheringJobs  []Job      `json:"gatheringJobs"`
	ProcessingJobs []*Recipes `json:"processingJobs"`
}

const (
	WoodcuttingJobType JobType = "woodcutting"
	MiningJobType      JobType = "mining"
	HarvestingJobType  JobType = "harvesting"
	SmeltingJobType    JobType = "smelting"
	WoodWorkingJobType JobType = "woodWorking"
)

type Reward struct {
	Items      []ItemWithQuantity `json:"items"`
	Experience int                `json:"experience"`
}

func (c *JobContainer) GetGatheringJobDefinition(jobDefID string) *Job {
	for _, def := range c.GatheringJobs {
		if def.ID == jobDefID {
			return &def
		}
	}

	return nil
}

func (c *JobContainer) GetProcessingJobDefinition(jobDefID string) *Recipes {
	for _, def := range c.ProcessingJobs {
		if def.ID == jobDefID {
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

func (j *Job) GetAffinty(element MonsterElement) float64 {
	if affinity, ok := j.Affinity.Elements[element]; ok {
		return affinity
	}
	return 1
}
