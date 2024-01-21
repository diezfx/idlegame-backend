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

func (s *JobService) StartSmeltingJob(ctx context.Context, userID, monsterID int, jobDefID string) (int, error) {
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

	taskDefinition := s.jobContainer.GetSmeltingJobDefinition(jobDefID)
	if taskDefinition == nil {
		return -1, fmt.Errorf("get job definition %s: %w", jobDefID, service.ErrJobTypeNotFound)
	}

	if taskDefinition.LevelRequirement > mon.Level() {
		return -1, service.ErrLevelRequirementNotMet
	}

	// start

	id, err := s.jobStorage.StoreNewJob(ctx, SmeltingJobType.String(), userID, monsterID, jobDefID)
	if err != nil {
		return -1, err
	}
	return id, nil
}

type Ingredient struct {
	Item  item.Item
	Count int
}

type Recipes struct {
	RecipeID string
	JobDefinition
	Ingredients []Ingredient
	OutputItem  item.Item
}

var smeltingJobs = []*Recipes{
	{
		JobDefinition: JobDefinition{
			JobDefID:         string(item.StoneBarType),
			JobType:          SmeltingJobType,
			LevelRequirement: 1,
			Duration:         time.Second * 3,
			Rewards: Reward{
				Items: []inventory.Item{
					{ItemDefID: string(item.StoneBarType), Quantity: 2},
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
