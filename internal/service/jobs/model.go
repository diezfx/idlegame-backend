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

type OreType string

const (
	StoneOreType   OreType = "Stone"
	CopperOreType  OreType = "Copper"
	IronOreType    OreType = "Iron"
	GoldOreType    OreType = "Gold"
	DiamondOreType OreType = "Diamond"
)

type CropType string

const (
	//low quality crop
	WheatCropType  CropType = "Wheat"
	CarrotCropType CropType = "Carrot"
	PotatoCropType CropType = "Potato"
)

func (t TreeType) String() string {
	return string(t)
}

func (t OreType) String() string {
	return string(t)
}

func (t CropType) String() string {
	return string(t)
}

type WoodCuttingJob struct {
	Job
	TreeType TreeType `json:"treeType"`
}

type MiningJob struct {
	Job
	OreType `json:"oreType"`
}

type HarvestingJob struct {
	Job
	CropType `json:"cropType"`
}

type JobContainer struct {
	woodcuttingDefs []WoodCuttingJobDefinition
	miningDefs      []MiningJobDefinition
	harvestingDefs  []HarvestingJobDefinition
}

// what kind of wood
type JobType string

const (
	WoodCuttingJobType JobType = "WoodCutting"
	MiningJobType      JobType = "Mining"
	HarvestingJobType  JobType = "Harvesting"
)

type Job struct {
	ID        int
	UserID    int
	StartedAt time.Time
	UpdatedAt time.Time
	Monsters  []int
	JobType   string
}

func FromWoodcuttingJob(j *storage.GatheringJob) *WoodCuttingJob {
	return &WoodCuttingJob{
		Job:      fromJob(j.Job),
		TreeType: TreeType(j.GatheringType),
	}
}

func FromMiningJob(j *storage.GatheringJob) *MiningJob {
	return &MiningJob{
		Job:     fromJob(j.Job),
		OreType: OreType(j.GatheringType),
	}
}

func FromHarvestingJob(j *storage.GatheringJob) *HarvestingJob {
	return &HarvestingJob{
		Job:      fromJob(j.Job),
		CropType: CropType(j.GatheringType),
	}
}

type Reward struct {
	Items []inventory.Item
	Exp   int
}

type JobDefinition struct {
	JobType          JobType
	LevelRequirement int
	Duration         time.Duration
	Rewards          Reward
}
type WoodCuttingJobDefinition struct {
	JobDefinition
	TreeType TreeType
}

type MiningJobDefinition struct {
	JobDefinition
	OreType OreType
}

type HarvestingJobDefinition struct {
	JobDefinition
	CropType CropType
}

func (c *JobContainer) GetWoodCuttingDefinition(treeType TreeType) *WoodCuttingJobDefinition {
	for _, def := range c.woodcuttingDefs {
		if def.TreeType == treeType {
			return &def
		}
	}
	return nil
}

func (c *JobContainer) GetMiningDefinition(oreType OreType) *MiningJobDefinition {
	for _, def := range c.miningDefs {
		if def.OreType == oreType {
			return &def
		}
	}
	return nil
}

func (c *JobContainer) GetHarvestingDefinition(cropType CropType) *HarvestingJobDefinition {
	for _, def := range c.harvestingDefs {
		if def.CropType == cropType {
			return &def
		}
	}
	return nil
}
