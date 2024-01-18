package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/service/monster"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

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
	LevelRequirement int
	Duration         time.Duration
	Rewards          Reward
}

var woodcuttingJobs = []WoodCuttingJobDefinition{
	{
		JobType:          SpruceType,
		LevelRequirement: 1,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: SpruceType.String(), Quantity: 1},
			},
			Exp: 1,
		},
	},
	{
		JobType:          BirchType,
		LevelRequirement: 2,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: BirchType.String(), Quantity: 1},
			},
			Exp: 2,
		},
	},
	{
		JobType:          PineType,
		LevelRequirement: 3,
		Duration:         time.Second * 3,
		Rewards: Reward{
			Items: []inventory.Item{
				{ItemDefID: PineType.String(), Quantity: 1},
			},
			Exp: 3,
		},
	},
}

func InitWoodCutting(itemContainer *item.ItemContainer) *WoodCuttingJobContainer {
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: SpruceType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: BirchType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: PineType.String(), Tags: []string{}})

	return &WoodCuttingJobContainer{
		defs: woodcuttingJobs,
	}
}

func (s *JobService) StartWoodCuttingJob(ctx context.Context, userID, monsterID int, treeType TreeType) (int, error) {
	// check if monster is not occupied
	_, err := s.jobStorage.GetJobByMonster(ctx, monsterID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return -1, fmt.Errorf("get job entry for %d: %w", monsterID, err)
	}
	if err == nil {
		return -1, service.ErrAlreadyStartedJob
	}

	// check if requirements are met

	storeMon, err := s.monsterStorage.GetMonsterByID(ctx, monsterID)
	if err != nil {
		return -1, fmt.Errorf("get monster information for monsterID %d: %w", monsterID, err)
	}
	mon := monster.MonsterFromStorage(storeMon)

	taskDefinition := s.woodContainer.GetDefinition(treeType)
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %d: %w", monsterID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewWoodCuttingJob(ctx, userID, monsterID, treeType.String())
	if err != nil {
		return -1, err
	}
	return id, nil
}

//getJob

func (s *JobService) GetWoodcuttingJob(ctx context.Context, id int) (*WoodCuttingJob, error) {
	job, err := s.jobStorage.GetWoodcuttingJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	return FromWoodcuttingJob(job), nil
}

func (s *JobService) StopWoodCuttingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("get job entry for jobID %d: %w", job.ID, err)
	}
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}

	// remove job
	err = s.jobStorage.DeleteWoodCuttingJob(ctx, id)
	if err != nil {
		return fmt.Errorf("delete job entry for jobID %d: %w", id, err)
	}
	return nil
}

func (s *JobService) UpdateWoodcuttingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetWoodcuttingJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.woodContainer.GetDefinition(job.TreeType)
	diff := job.UpdatedAt.Sub(job.StartedAt)
	steps := diff / jobDefintion.Duration // the ticks from the beginning

	executionCount := 0
	for nextTick := job.StartedAt.Add(jobDefintion.Duration * steps); nextTick.Before(now); nextTick = nextTick.Add(jobDefintion.Duration) {
		if nextTick.After(job.UpdatedAt) {
			executionCount++
		}
	}
	if executionCount == 0 {
		return nil
	}

	rewards := calculateRewards(jobDefintion.Rewards, executionCount)

	err = s.inventoryStorage.AddItems(ctx, toInventoryEntries(job.UserID, rewards.Items))
	if err != nil {
		return fmt.Errorf("add items for userID %d: %w", job.UserID, err)
	}

	_, err = s.monsterStorage.AddMonsterExperience(ctx, job.Monsters[0], rewards.Exp)
	if err != nil {
		return fmt.Errorf("add exp for userID %d: %w", job.UserID, err)
	}

	err = s.jobStorage.UpdateJobUpdatedAt(ctx, id, now)
	if err != nil {
		return fmt.Errorf("update job entry for jobID %d: %w", id, err)
	}
	return nil

}

func toInventoryEntries(userId int, item []inventory.Item) []storage.InventoryEntry {
	entries := []storage.InventoryEntry{}
	for _, i := range item {
		entries = append(entries, storage.InventoryEntry{
			UserID:    userId,
			ItemDefID: i.ItemDefID,
			Quantity:  i.Quantity,
		})
	}
	return entries
}

func calculateRewards(rewards Reward, executionCount int) Reward {

	var rewardItems = []inventory.Item{}
	for _, item := range rewards.Items {
		rewardItems = append(rewardItems, inventory.Item{
			Quantity:  item.Quantity * executionCount,
			ItemDefID: item.ItemDefID,
		})
	}
	return Reward{
		Items: rewardItems,
		Exp:   rewards.Exp * executionCount,
	}
}
