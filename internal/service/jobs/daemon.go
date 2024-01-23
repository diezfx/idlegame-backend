package jobs

import (
	"context"
	"time"

	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

const tickTime = time.Second

// Daemon queries for jobs in the background and takes care of the effects.
type Daemon struct {
	jobService DaemonJobService
}

type DaemonJobService interface {
	GetJobs(ctx context.Context) ([]Job, error)
	UpdateGatheringJob(ctx context.Context, id int) error
	UpdateProcessingJob(ctx context.Context, id int) error
}

func NewDaemon(jobService DaemonJobService) *Daemon {
	return &Daemon{
		jobService: jobService,
	}
}

// TODO use cronlibrary
// TODO err chan
func (d *Daemon) Run(ctx context.Context) error {
	ticker := time.NewTicker(tickTime)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			jobs, err := d.jobService.GetJobs(ctx)
			if err != nil {
				logger.Error(ctx, err).Msg("get jobs")
				return err
			}
			for _, job := range jobs {
				if job.JobType == masterdata.WoodcuttingJobType.String() || job.JobType == masterdata.MiningJobType.String() || job.JobType == masterdata.HarvestingJobType.String() {
					err = d.jobService.UpdateGatheringJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update woodcutting job")
						return err
					}
				}
				if job.JobType == masterdata.SmeltingJobType.String() || job.JobType == masterdata.WoodWorkingJobType.String() {
					err = d.jobService.UpdateProcessingJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update smelting job")
						return err
					}
				}
			}
		}
	}
}
