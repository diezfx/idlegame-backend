package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobService struct {
	jobStorage       JobStorage
	monsterStorage   MonsterStorage
	inventoryStorage InventoryStorage
	woodContainer    WoodCuttingJobContainer
}

func New(jobStorage JobStorage, monsterStorage MonsterStorage, inventoryStorage InventoryStorage, itemContainer *item.ItemContainer) *JobService {
	return &JobService{
		jobStorage: jobStorage, monsterStorage: monsterStorage, inventoryStorage: inventoryStorage,
		woodContainer: *InitWoodCutting(itemContainer),
	}
}

func (s *JobService) GetJob(ctx context.Context, id int) (*Job, error) {
	// check if job exists
	storageJob, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", storageJob.ID, err)
	}
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	j := fromJob(*storageJob)
	return &j, nil
}

func (t JobType) String() string {
	return string(t)
}

func (s *JobService) GetJobs(ctx context.Context) ([]Job, error) {
	storageJobs, err := s.jobStorage.GetJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get jobs: %w", err)
	}

	jobs := make([]Job, 0, len(storageJobs))
	for _, storageJob := range storageJobs {
		jobs = append(jobs, fromJob(storageJob))
	}
	return jobs, nil
}

func fromJob(j storage.Job) Job {
	return Job{
		ID:        j.ID,
		StartedAt: j.StartedAt,
		UpdatedAt: j.UpdatedAt,
		Monsters:  j.Monsters,
		JobType:   j.JobType,
		UserID:    j.UserID,
	}
}
