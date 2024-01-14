package jobs

import (
	"context"
	"time"

	"github.com/diezfx/idlegame-backend/pkg/logger"
)

const tickTime = time.Second

// Daemon queries for jobs in the background and takes care of the effects.
type Daemon struct {
	jobService DaemonJobService
}

type DaemonJobService interface {
	GetJobs(ctx context.Context) ([]Job, error)
	UpdateWoodcuttingJob(ctx context.Context, id int) error
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
			logger.Debug(ctx).Msg("daemon tick for jobs")
			jobs, err := d.jobService.GetJobs(ctx)
			if err != nil {
				logger.Error(ctx, err).Msg("get jobs")
				return err
			}
			logger.Debug(ctx).Any("jobs", jobs).Msg("jobs")
			for _, job := range jobs {
				if job.JobType == WoodCuttingJobType.String() {
					err = d.jobService.UpdateWoodcuttingJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update woodcutting job")
						return err
					}
				}

			}
		}

	}

}
