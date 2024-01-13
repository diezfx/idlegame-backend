package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobService struct {
	jobStorage     JobStorage
	monsterStorage MonsterStorage
	woodContainer  WoodCuttingJobContainer
}

func New(jobStorage JobStorage, monsterStorage MonsterStorage, itemContainer *item.ItemContainer) *JobService {
	return &JobService{
		jobStorage: jobStorage, monsterStorage: monsterStorage,
		woodContainer: *InitWoodCutting(itemContainer),
	}
}

func (s *JobService) GetJob(ctx context.Context, id int) (*Job, error) {
	// check if job exists
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", job.ID, err)
	}
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	return &Job{
		ID:        job.ID,
		StartedAt: job.StartedAt,
		Monsters:  job.Monsters,
		JobType:   job.JobType,
	}, nil
}

func (t JobType) String() string {
	return string(t)
}
