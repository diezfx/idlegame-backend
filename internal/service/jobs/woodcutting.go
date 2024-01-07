package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

// what kind of wood
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
	ID       int
	Monster  int
	TreeType TreeType
}

// different tree types have
// level requirements, durations, exp gains

type WoodCuttingJobContainer struct {
	defs []WoodCuttingJobDefinition
}

func (c *WoodCuttingJobContainer) GetDefinition(treeType TreeType) *WoodCuttingJobDefinition {
	for _, def := range c.defs {
		if def.JobType == treeType {
			return &def
		}
	}
	return nil
}

type WoodCuttingJobDefinition struct {
	JobType          TreeType
	ExpGain          int
	LevelRequirement int
	Duration         time.Duration
}

var woodcuttingJobs = []WoodCuttingJobDefinition{
	{
		ExpGain:          1,
		JobType:          SpruceType,
		LevelRequirement: 1,
		Duration:         time.Second,
	},
	{
		ExpGain:          2,
		JobType:          BirchType,
		LevelRequirement: 2,
		Duration:         time.Second * 3,
	},
	{
		ExpGain:          3,
		JobType:          PineType,
		LevelRequirement: 3,
		Duration:         time.Second * 3,
	},
}

func InitWoodCutting(itemContainer *item.ItemContainer) {
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: SpruceType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: BirchType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: PineType.String(), Tags: []string{}})
}

func (s *JobService) StartWoodCuttingJob(ctx context.Context, job WoodCuttingJob) error {
	// check if monster is not occupied
	_, err := s.jobStorage.GetMonsterEntry(ctx, job.Monster)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("get job entry for %s: %w", job.Monster, err)
	}
	if err == nil {
		return service.ErrAlreadyStartedJob
	}

	// check if requirements are meant

	mon, err := s.monsterStorage.GetMonster(job.Monster)
	if err != nil {
		return fmt.Errorf("get monster information for %s: %w", job.Monster, err)
	}

	taskDefinition := s.woodContainer.GetDefinition(job.TreeType)
	if taskDefinition == nil {
		return fmt.Errorf("get job definition %s: %w", job.Monster, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return service.ErrJobTypeNotFound
	}

	// start

	s.jobStorage.StoreMonsterEntry()

	// add entry into monsterJobEntry, add job into job table
}
