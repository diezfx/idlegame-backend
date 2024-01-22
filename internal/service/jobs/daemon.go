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
	UpdateWoodcuttingJob(ctx context.Context, id int) error
	UpdateMiningJob(ctx context.Context, id int) error
	UpdateHarvestingJob(ctx context.Context, id int) error
	UpdateSmeltingJob(ctx context.Context, id int) error
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
				if job.JobType == masterdata.WoodCuttingJobType.String() {
					err = d.jobService.UpdateWoodcuttingJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update woodcutting job")
						return err
					}
				}
				if job.JobType == masterdata.MiningJobType.String() {
					err = d.jobService.UpdateMiningJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update mining job")
						return err
					}
				}
				if job.JobType == masterdata.HarvestingJobType.String() {
					err = d.jobService.UpdateHarvestingJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update harvesting job")
						return err
					}
				}
				if job.JobType == masterdata.SmeltingJobType.String() {
					err = d.jobService.UpdateSmeltingJob(ctx, job.ID)
					if err != nil {
						logger.Error(ctx, err).Msg("update smelting job")
						return err
					}
				}
			}
		}
	}
}
