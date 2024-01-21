package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/logger"
)

type JobService struct {
	jobStorage       JobStorage
	monsterStorage   MonsterStorage
	inventoryStorage InventoryStorage
	woodContainer    JobContainer
}

func New(jobStorage JobStorage, monsterStorage MonsterStorage, inventoryStorage InventoryStorage, itemContainer *item.ItemContainer) *JobService {
	return &JobService{
		jobStorage: jobStorage, monsterStorage: monsterStorage, inventoryStorage: inventoryStorage,
		woodContainer: *InitJobs(itemContainer),
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

func InitJobs(itemContainer *item.ItemContainer) *JobContainer {
	// Woodcutting
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: SpruceType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: BirchType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: PineType.String(), Tags: []string{}})

	// Mining
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: StoneOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: CopperOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: IronOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: GoldOreType.String(), Tags: []string{}})
	itemContainer.AddItemDefinition(item.ItemDefinition{ID: DiamondOreType.String(), Tags: []string{}})

	return &JobContainer{
		woodcuttingDefs: woodcuttingJobs,
		miningDefs:      miningJobs,
		harvestingDefs:  harvestingJobs,
	}
}

func (s *JobService) StopJob(ctx context.Context, id int) error {
	// check if job exists
	job, err := s.jobStorage.GetJobByID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("get job entry for jobID %d: %w", job.ID, err)
	}
	if err != nil {
		return fmt.Errorf("get job entry for jobID %d: %w", id, err)
	}

	if JobType(job.JobType) == WoodCuttingJobType || JobType(job.JobType) == MiningJobType || JobType(job.JobType) == HarvestingJobType {
		// remove job
		err = s.jobStorage.DeleteGatheringJob(ctx, id)
		if err != nil {
			return fmt.Errorf("delete job entry for jobID %d: %w", id, err)
		}
		return nil
	} else {
		logger.Fatal(ctx, nil).String("jobType", job.JobType).Msg("not implemented yet")
	}
	return nil
}
