//nolint:dupl // fine for now
package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

func (s *JobService) StartGatheringJob(ctx context.Context, userID, monsterID int, jobDefID string) (int, error) {

	taskDefinition := s.masterdata.Jobs.GetGatheringJobDefinition(jobDefID)
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %d: %w", monsterID, service.ErrJobTypeNotFound)
	}
	// check if monster is not occupied
	_, err := s.jobStorage.GetJobByMonster(ctx, monsterID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return -1, fmt.Errorf("get job entry for %d: %w", monsterID, err)
	}
	if err == nil {
		return -1, service.ErrAlreadyStartedJob
	}

	// check if requirements are met

	mon, err := s.monsterStorage.GetMonsterByID(ctx, monsterID)
	if err != nil {
		return -1, fmt.Errorf("get monster information for monsterID %d: %w", monsterID, err)
	}
	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, taskDefinition.JobType.String(), userID, monsterID, jobDefID)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *JobService) UpdateGatheringJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.masterdata.Jobs.GetGatheringJobDefinition(job.JobDefID)

	monster, err := s.monsterStorage.GetMonsterByID(ctx, job.Monsters[0])
	if err != nil {
		return fmt.Errorf("get monster information for monsterID %d: %w", job.Monsters[0], err)
	}
	duration := jobDefintion.Duration.Duration() / time.Duration(jobDefintion.GetAffinty(monster.Element()))

	executionCount := calculateTicks(job, duration, now)

	rewards := calculateRewards(jobDefintion.Rewards, executionCount)

	err = s.inventoryStorage.AddItems(ctx, toInventoryEntries(job.UserID, rewards.Items))
	if err != nil {
		return fmt.Errorf("add items for userID %d: %w", job.UserID, err)
	}

	_, err = s.monsterStorage.AddMonsterExperience(ctx, job.Monsters[0], rewards.Experience)
	if err != nil {
		return fmt.Errorf("add exp for userID %d: %w", job.UserID, err)
	}

	err = s.jobStorage.UpdateJobUpdatedAt(ctx, id, now)
	if err != nil {
		return fmt.Errorf("update job entry for jobID %d: %w", id, err)
	}
	return nil
}

func toInventoryEntries(userID int, itm []masterdata.ItemWithQuantity) []storage.InventoryEntry {
	entries := []storage.InventoryEntry{}
	for _, i := range itm {
		entries = append(entries, storage.InventoryEntry{
			UserID:    userID,
			ItemDefID: i.ID,
			Quantity:  i.Quantity,
		})
	}
	return entries
}

func costToInventoryEntries(userID int, itm []masterdata.ItemWithQuantity) []storage.InventoryEntry {
	entries := []storage.InventoryEntry{}
	for _, i := range itm {
		entries = append(entries, storage.InventoryEntry{
			UserID:    userID,
			ItemDefID: i.ID,
			Quantity:  -i.Quantity,
		})
	}
	return entries
}

func calculateRewards(rewards masterdata.Reward, executionCount int) masterdata.Reward {
	var rewardItems = []masterdata.ItemWithQuantity{}
	for _, item := range rewards.Items {
		rewardItems = append(rewardItems, masterdata.ItemWithQuantity{
			Quantity: item.Quantity * executionCount,
			ID:       item.ID,
		})
	}
	return masterdata.Reward{
		Items:      rewardItems,
		Experience: rewards.Experience * executionCount,
	}
}

func calculateCosts(costs []masterdata.ItemWithQuantity, executionCount int) []masterdata.ItemWithQuantity {
	var costItems = []masterdata.ItemWithQuantity{}
	for _, item := range costs {
		costItems = append(costItems, masterdata.ItemWithQuantity{
			ID:       item.ID,
			Quantity: item.Quantity * executionCount,
		})
	}
	return costItems
}
