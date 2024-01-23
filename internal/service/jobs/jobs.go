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

type JobService struct {
	jobStorage       JobStorage
	monsterStorage   MonsterStorage
	inventoryStorage InventoryStorage
	masterdata       *masterdata.Container
}

func New(jobStorage JobStorage, monsterStorage MonsterStorage, inventoryStorage InventoryStorage, masterdata *masterdata.Container) *JobService {
	return &JobService{
		jobStorage:       jobStorage,
		monsterStorage:   monsterStorage,
		inventoryStorage: inventoryStorage,
		masterdata:       masterdata,
	}
}

func (s *JobService) GetJob(ctx context.Context, id int) (*Job, error) {
	// check if job exists
	storageJob, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		return nil, service.ErrJobNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}
	j := fromJob(*storageJob)
	return &j, nil
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
		JobDefID:  j.JobDefID,
		StartedAt: j.StartedAt,
		UpdatedAt: j.UpdatedAt,
		Monsters:  j.Monsters,
		JobType:   j.JobType,
		UserID:    j.UserID,
	}
}

func calculateTicks(job Job, duration time.Duration, now time.Time) int {
	diff := job.UpdatedAt.Sub(job.StartedAt)
	steps := diff / duration // the ticks from the beginning
	executionCount := 0
	for nextTick := job.StartedAt.Add(duration * steps); nextTick.Before(now); nextTick = nextTick.Add(duration) {
		if nextTick.After(job.UpdatedAt) {
			executionCount++
		}
	}
	return executionCount
}

func (s *JobService) StopJob(ctx context.Context, id int) error {
	// check if job exists
	_, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return service.ErrJobNotFound
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
