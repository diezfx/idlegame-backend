package jobs

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/diezfx/idlegame-backend/internal/service"
	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

func (s *JobService) StartProcessingJob(ctx context.Context, userID, monsterID int, jobDefID string) (int, error) {
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

	taskDefinition := s.masterdata.Jobs.GetProcessingJobDefinition(jobDefID)
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %s: %w", jobDefID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	inventoryStr, err := s.inventoryStorage.GetInventory(ctx, userID)
	if err != nil {
		return -1, fmt.Errorf("get inventory for userID %d: %w", userID, err)
	}
	inv := inventory.ToInventoryFromStorageEntries(inventoryStr, userID)
	maxRuns := calculateMaxRuns(inv, taskDefinition)
	if maxRuns == 0 {
		return -1, service.ErrNotEnoughItems
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, taskDefinition.JobType.String(), userID, monsterID, jobDefID)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *JobService) UpdateProcessingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.masterdata.Jobs.GetProcessingJobDefinition(job.JobDefID)

	executionCount := calculateTicks(job, jobDefintion.Duration.Duration(), now)

	inventoryStr, err := s.inventoryStorage.GetInventory(ctx, job.UserID)
	if err != nil {
		return fmt.Errorf("get inventory for userID %d: %w", job.UserID, err)
	}
	inv := inventory.ToInventoryFromStorageEntries(inventoryStr, job.UserID)

	// check that enough items are there
	maxRuns := calculateMaxRuns(inv, jobDefintion)

	if maxRuns < executionCount {
		logger.Debug(ctx).Int("jobID", id).Msg("drop job")
		err = s.jobStorage.DeleteJobByID(ctx, id)
		if err != nil {
			return fmt.Errorf("delete job entry for jobID %d: %w", id, err)
		}
		executionCount = maxRuns
	}

	if maxRuns == 0 {
		return nil
	}

	rewards := calculateRewards(jobDefintion.Rewards, executionCount)
	// item to get
	costs := calculateCosts(jobDefintion.Ingredients, executionCount)

	err = s.inventoryStorage.AddItems(ctx, costToInventoryEntries(job.UserID, costs))
	if err != nil {
		return fmt.Errorf("remove items for userID %d: %w", job.UserID, err)
	}

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

func calculateMaxRuns(inv *inventory.Inventory, jobDef *masterdata.Recipes) int {
	maxRuns := math.MaxInt32
	for _, cost := range jobDef.Ingredients {
		found := false
		for _, item := range inv.Items {
			if item.ID == cost.ID {
				max := item.Quantity / cost.Quantity
				if max < maxRuns {
					maxRuns = max
					found = true
				}
			}
		}
		if !found {
			return 0
		}
	}
	return maxRuns
}
