package masterdata

import (
	"time"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/item"
)

type Job struct {
	ID               string        `json:"id"`
	JobType          JobType       `json:"jobType"`
	LevelRequirement int           `json:"levelRequirement"`
	Duration         time.Duration `json:"duration"`
	Rewards          Reward        `json:"rewards"`
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
	WoodCuttingJobType JobType = "WoodCutting"
	MiningJobType      JobType = "Mining"
	HarvestingJobType  JobType = "Harvesting"
	SmeltingJobType    JobType = "Smelting"
)

type Reward struct {
	Items []inventory.Item
	Exp   int
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

var woodcuttingJobs = []Job{

	{
		ID:               item.SpruceType.String(),
		JobType:          WoodCuttingJobType,
		LevelRequirement: 1,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.SpruceType.String(), Quantity: 1},
			},
			Exp: 1,
		},
	},
	{
		ID:               item.BirchType.String(),
		JobType:          WoodCuttingJobType,
		LevelRequirement: 2,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.BirchType.String(), Quantity: 1},
			},
			Exp: 2,
		},
	},

	{
		ID:               item.PineType.String(),
		JobType:          WoodCuttingJobType,
		LevelRequirement: 3,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.PineType.String(), Quantity: 1},
			},
			Exp: 3,
		},
	},
}

type Ingredient struct {
	Item  item.Item
	Count int
}

type Recipes struct {
	ID string
	Job
	Ingredients []Ingredient
	OutputItem  item.Item
}

var smeltingJobs = []*Recipes{
	{
		Job: Job{
			ID:               string(item.StoneBarType),
			JobType:          SmeltingJobType,
			LevelRequirement: 1,
			Duration:         time.Second * 3,
			Rewards: Reward{
				Items: []inventory.Item{
					{ItemDefID: string(item.StoneBarType), Quantity: 1},
				},
				Exp: 1,
			},
		},
		Ingredients: []Ingredient{
			{
				Item:  item.Item(item.StoneOreType),
				Count: 2,
			},
		},
	},
}

var miningJobs = []Job{
	{
		ID:               item.StoneOreType.String(),
		JobType:          MiningJobType,
		LevelRequirement: 1,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.StoneOreType.String(), Quantity: 1},
			},
			Exp: 1,
		},
	},
	{
		ID:               item.CopperOreType.String(),
		JobType:          MiningJobType,
		LevelRequirement: 2,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.CopperOreType.String(), Quantity: 1},
			},
			Exp: 2,
		},
	},
	{
		ID:               item.IronOreType.String(),
		JobType:          MiningJobType,
		LevelRequirement: 3,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.IronOreType.String(), Quantity: 1},
			},
			Exp: 3,
		},
	},
	{
		ID:               item.GoldOreType.String(),
		JobType:          MiningJobType,
		LevelRequirement: 4,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.GoldOreType.String(), Quantity: 1},
			},
			Exp: 4,
		},
	},
	{
		ID:               item.DiamondOreType.String(),
		JobType:          MiningJobType,
		LevelRequirement: 5,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.DiamondOreType.String(), Quantity: 1},
			},
			Exp: 5,
		},
	},
}

func InitJobs(itemContainer *item.ItemContainer) *JobContainer {
	// Woodcutting
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.SpruceType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.BirchType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.PineType.String(), Tags: []string{}})

	// Mining
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.StoneOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.CopperOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.IronOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.GoldOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.DiamondOreType.String(), Tags: []string{}})

	// Harvesting
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.CarrotCropType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.PotatoCropType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: item.WheatCropType.String(), Tags: []string{}})

	gatheringJobs := []Job{}
	gatheringJobs = append(gatheringJobs, miningJobs...)
	gatheringJobs = append(gatheringJobs, woodcuttingJobs...)
	gatheringJobs = append(gatheringJobs, harvestingJobs...)
	return &JobContainer{
		GatheringJobs:  gatheringJobs,
		ProcessingJobs: smeltingJobs,
	}
}

var harvestingJobs = []Job{

	{
		ID:               item.WheatCropType.String(),
		JobType:          HarvestingJobType,
		LevelRequirement: 1,
		Duration:         time.Second * 5,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.WheatCropType.String(), Quantity: 1},
			},
			Exp: 1,
		},
	},
	{
		ID:               item.CarrotCropType.String(),
		JobType:          HarvestingJobType,
		LevelRequirement: 2,
		Duration:         time.Second * 5,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.CarrotCropType.String(), Quantity: 1},
			},
			Exp: 2,
		},
	},
	{
		ID:               item.PotatoCropType.String(),
		JobType:          HarvestingJobType,
		LevelRequirement: 3,
		Duration:         time.Second * 5,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: item.PotatoCropType.String(), Quantity: 1},
			},
			Exp: 3,
		},
	},
}
