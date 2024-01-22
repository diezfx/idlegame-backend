//nolint:dupl // fine for now
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
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

func (s *JobService) StartWoodCuttingJob(ctx context.Context, userID, monsterID int, treeType item.TreeType) (int, error) {
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

	taskDefinition := s.jobContainer.GetGatheringJobDefinition(masterdata.WoodCuttingJobType, treeType.String())
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %d: %w", monsterID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, masterdata.WoodCuttingJobType.String(), userID, monsterID, treeType.String())
	if err != nil {
		return -1, err
	}
	return id, nil
}

// getJob

func (s *JobService) GetWoodcuttingJob(ctx context.Context, id int) (*WoodCuttingJob, error) {
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	return FromWoodcuttingJob(job), nil
}

func (s *JobService) UpdateWoodcuttingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetWoodcuttingJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.jobContainer.GetGatheringJobDefinition(masterdata.WoodCuttingJobType, job.TreeType.String())

	executionCount := calculateTicks(job.Job, jobDefintion.Duration, now)

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

func toInventoryEntries(userID int, itm []inventory.Item) []storage.InventoryEntry {
	entries := []storage.InventoryEntry{}
	for _, i := range itm {
		entries = append(entries, storage.InventoryEntry{
			UserID:    userID,
			ItemDefID: i.ItemDefID,
			Quantity:  i.Quantity,
		})
	}
	return entries
}

func costToInventoryEntries(userID int, itm []masterdata.Ingredient) []storage.InventoryEntry {
	entries := []storage.InventoryEntry{}
	for _, i := range itm {
		entries = append(entries, storage.InventoryEntry{
			UserID:    userID,
			ItemDefID: string(i.Item),
			Quantity:  -i.Count,
		})
	}
	return entries
}

func calculateRewards(rewards masterdata.Reward, executionCount int) masterdata.Reward {
	var rewardItems = []inventory.Item{}
	for _, item := range rewards.Items {
		rewardItems = append(rewardItems, inventory.Item{
			Quantity:  item.Quantity * executionCount,
			ItemDefID: item.ItemDefID,
		})
	}
	return masterdata.Reward{
		Items: rewardItems,
		Exp:   rewards.Exp * executionCount,
	}
}

func calculateCosts(costs []masterdata.Ingredient, executionCount int) []masterdata.Ingredient {
	var costItems = []masterdata.Ingredient{}
	for _, item := range costs {
		costItems = append(costItems, masterdata.Ingredient{
			Item:  item.Item,
			Count: item.Count * executionCount,
		})
	}
	return costItems
}
