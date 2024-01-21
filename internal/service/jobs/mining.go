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

var miningJobs = []JobDefinition{
	JobDefinition{
		JobDefID:         item.StoneOreType.String(),
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
	JobDefinition{
		JobDefID:         item.CopperOreType.String(),
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
	JobDefinition{
		JobDefID:         item.IronOreType.String(),
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
	JobDefinition{
		JobDefID:         item.GoldOreType.String(),
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
	JobDefinition{
		JobDefID:         item.DiamondOreType.String(),
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

func (s *JobService) StartMiningJob(ctx context.Context, userID, monsterID int, oreType item.OreType) (int, error) {
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

	taskDefinition := s.jobContainer.GetGatheringJobDefinition(MiningJobType, oreType.String())
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %d: %w", monsterID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, MiningJobType.String(), userID, monsterID, oreType.String())
	if err != nil {
		return -1, err
	}
	return id, nil
}

//getJob

func (s *JobService) GetMiningJob(ctx context.Context, id int) (*MiningJob, error) {
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	return FromMiningJob(job), nil
}

func (s *JobService) UpdateMiningJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.GetMiningJob(ctx, id)
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	now := time.Now()
	jobDefintion := s.jobContainer.GetGatheringJobDefinition(MiningJobType, string(job.OreType))
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
