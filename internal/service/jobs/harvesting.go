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

var harvestingJobs = []JobDefinition{

	JobDefinition{
		JobDefID:         item.WheatCropType.String(),
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
	JobDefinition{
		JobDefID:         item.CarrotCropType.String(),
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
	JobDefinition{
		JobDefID:         item.PotatoCropType.String(),
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

func (s *JobService) StartHarvestingJob(ctx context.Context, userID, monsterID int, cropType item.CropType) (int, error) {
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
		return -1, fmt.Errorf("get monster information for monsterID %d: %w", cropType, err)
	}
	mon := monster.MonsterFromStorage(storeMon)

	taskDefinition := s.jobContainer.GetGatheringJobDefinition(HarvestingJobType, cropType.String())
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %d: %w", monsterID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, HarvestingJobType.String(), userID, monsterID, cropType.String())
	if err != nil {
		return -1, err
	}
	return id, nil
}

//getJob

func (s *JobService) GetHarvestingJob(ctx context.Context, id int) (*HarvestingJob, error) {
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	return FromHarvestingJob(job), nil
}

func (s *JobService) StopHarvestingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("get job entry for jobID %d: %w", job.ID, err)
	}
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}

	// remove job
	err = s.jobStorage.DeleteJobByID(ctx, id)
	if err != nil {
		return fmt.Errorf("delete job entry for jobID %d: %w", id, err)
	}
	return nil
}

func (s *JobService) UpdateHarvestingJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetHarvestingJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.jobContainer.GetGatheringJobDefinition(HarvestingJobType, job.CropType.String())
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
