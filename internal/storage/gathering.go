package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/pkg/db"
)

func (c *Client) GetGatheringJobByID(ctx context.Context, id int) (*GatheringJob, error) {
	const getGatheringJobByID = `
	SELECT j.id,j.user_id,j.started_at,j.updated_at, m.monster_id, w.gathering_type
	FROM jobs as j
	LEFT JOIN gathering_jobs as w
	ON j.id=w.job_id
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []struct {
		job
		jobMonster
		GatheringType string
	}
	err := c.dbClient.Select(ctx, &res, getGatheringJobByID, id)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	if len(res) > 1 {
		return nil, fmt.Errorf("multiple gathering jobs found")
	}
	gatheringJob := toGatheringJob(res)
	return &gatheringJob, nil
}

func (c *Client) StoreNewGatheringJob(ctx context.Context, jobType string, userID, monsterID int, gatheringType string) (int, error) {
	var jobID int

	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {

		err := tx.Get(ctx, &jobID, insertJobQuery, userID, time.Now(), jobType)
		if err != nil {
			return fmt.Errorf("insert job: %w", err)
		}
		_, err = tx.Exec(ctx, insertMonsterQuery, jobID, monsterID)
		if err != nil {
			return fmt.Errorf("insert job_monster: %w", err)
		}

		const insertGatheringJobQuery = `
		INSERT INTO gathering_jobs(job_id,gathering_type)
		values($1,$2)`
		_, err = tx.Exec(ctx, insertGatheringJobQuery, jobID, gatheringType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return -1, fmt.Errorf("transaction failed: %w", err)
	}
	return jobID, nil
}

func (c *Client) DeleteGatheringJob(ctx context.Context, jobID int) error {
	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		_, err := tx.Exec(ctx, deleteMonsterEntriesQuery, jobID)
		if err != nil {
			return err
		}

		const deleteGatheringJobQuery = `
		DELETE FROM gathering_jobs
		WHERE job_id=$1
		`
		_, err = tx.Exec(ctx, deleteGatheringJobQuery, jobID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, deleteJobQuery, jobID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func toGatheringJob(res []struct {
	job
	jobMonster
	GatheringType string
},
) GatheringJob {
	job := GatheringJob{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, entry.MonsterID)
	}
	job.ID = res[0].ID
	job.UserID = res[0].UserID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	job.UpdatedAt = res[0].UpdatedAt
	job.GatheringType = res[0].GatheringType

	return job
}
