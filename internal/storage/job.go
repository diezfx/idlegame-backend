package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/pkg/db"
)

const insertJobQuery = `
INSERT INTO jobs(job_def_id,user_id,started_at,updated_at,job_type)
values($1,$2,$3,$3,$4)
RETURNING id`

const insertMonsterQuery = `
INSERT INTO job_monsters (job_id, monster_id)
values($1,$2)`

const deleteMonsterEntriesQuery = `
DELETE FROM job_monsters
WHERE job_id=$1
`

const deleteJobQuery = `
DELETE FROM jobs
WHERE id=$1
`

func (c *Client) UpdateJobUpdatedAt(ctx context.Context, id int, updatedAt time.Time) error {
	const updateJobQuery = `
	UPDATE jobs
	SET updated_at=$1
	WHERE id=$2
	`
	_, err := c.dbClient.Exec(ctx, updateJobQuery, updatedAt, id)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}
	return nil
}

func (c *Client) GetJobByMonster(ctx context.Context, monID int) (*Job, error) {
	const getJobByMonsterQuery = `
	SELECT j.id,j.started_at,j.updated_at,j.job_type, m.monster_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE m.monster_id=$1`

	var res []getJobsQueryResult
	err := c.dbClient.Select(ctx, &res, getJobByMonsterQuery, monID)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}

	job := toJob(res)

	return &job, nil
}

func (c *Client) GetJobByID(ctx context.Context, id int) (*Job, error) {
	const getJobByID = `
	SELECT j.id,j.job_def_id,j.started_at,j.updated_at,j.job_type, m.monster_id, j.user_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []getJobsQueryResult
	err := c.dbClient.Select(ctx, &res, getJobByID, id)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}

	job := toJob(res)
	return &job, nil
}

func (c *Client) StoreNewJob(ctx context.Context, jobType string, userID, monsterID int, jobDefID string) (int, error) {
	var jobID int

	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		err := tx.Get(ctx, &jobID, insertJobQuery, jobDefID, userID, time.Now(), jobType)
		if err != nil {
			return fmt.Errorf("insert job: %w", err)
		}
		_, err = tx.Exec(ctx, insertMonsterQuery, jobID, monsterID)
		if err != nil {
			return fmt.Errorf("insert job_monster: %w", err)
		}
		return nil
	})
	if err != nil {
		return -1, fmt.Errorf("transaction failed: %w", err)
	}
	return jobID, nil
}

type getJobsQueryResult struct {
	job
	jobMonster
}

func (c *Client) GetJobs(ctx context.Context) ([]Job, error) {
	const getJobsQuery = `
	SELECT j.id,j.started_at,j.updated_at,j.job_type, m.monster_id, j.user_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	`
	var res []getJobsQueryResult
	err := c.dbClient.Select(ctx, &res, getJobsQuery)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}

	jobMap := make(map[int][]getJobsQueryResult)
	for _, entry := range res {
		jobMap[entry.ID] = append(jobMap[entry.ID], entry)
	}

	jobs := make([]Job, 0, len(res))
	for _, entry := range jobMap {
		job := toJob(entry)
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func toJob(res []getJobsQueryResult) Job {
	job := Job{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, entry.MonsterID)
	}
	job.ID = res[0].ID
	job.JobDefID = res[0].JobDefID
	job.UserID = res[0].UserID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	job.UpdatedAt = res[0].UpdatedAt
	return job
}

func (c *Client) DeleteJobByID(ctx context.Context, jobID int) error {
	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		_, err := tx.Exec(ctx, deleteMonsterEntriesQuery, jobID)
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
