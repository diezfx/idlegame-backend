package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/georgysavva/scany/sqlscan"
)

func (c *Client) StoreNewWoodCuttingJob(ctx context.Context, userID, monsterID int, woodType string) (int, error) {
	var jobID int

	err := c.withTx(ctx, func(tx *sql.Tx) error {
		const insertJobQuery = `
		INSERT INTO jobs(user_id,started_at,updated_at,job_type)
		values($1,$2,$2,$3)
		RETURNING id`

		err := tx.QueryRowContext(ctx, insertJobQuery, userID, time.Now(), "WoodCutting").Scan(&jobID)
		if err != nil {
			return fmt.Errorf("insert job: %w", err)
		}

		const insertMonsterQuery = `
		INSERT INTO job_monsters (job_id, monster_id)
		values($1,$2)`
		_, err = tx.ExecContext(ctx, insertMonsterQuery, jobID, monsterID)
		if err != nil {
			return fmt.Errorf("insert job_monster: %w", err)
		}

		const insertWoodQuery = `
		INSERT INTO woodcutting_jobs(job_id,tree_type)
		values($1,$2)`
		_, err = tx.ExecContext(ctx, insertWoodQuery, jobID, woodType)
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

func (c *Client) DeleteWoodCuttingJob(ctx context.Context, jobID int) error {
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		const deleteJobQuery = `
		DELETE FROM jobs
		WHERE id=$1
		`
		_, err := tx.ExecContext(ctx, deleteJobQuery, jobID)
		if err != nil {
			return err
		}

		const deleteMonsterEntriesQuery = `
		DELETE FROM job_monsters
		WHERE job_id=$1
		`
		_, err = tx.ExecContext(ctx, deleteMonsterEntriesQuery, jobID)
		if err != nil {
			return err
		}
		const deleteWoodCuttingJobQuery = `
		DELETE FROM woodcutting_jobs
		WHERE job_id=$1
		`
		_, err = tx.ExecContext(ctx, deleteWoodCuttingJobQuery, jobID)
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

func (c *Client) UpdateJobUpdatedAt(ctx context.Context, id int, updatedAt time.Time) error {
	const updateJobQuery = `
	UPDATE jobs
	SET updated_at=$1
	WHERE id=$2
	`
	_, err := c.conn.ExecContext(ctx, updateJobQuery, updatedAt, id)
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
	err := sqlscan.Select(ctx, c.conn.DB, &res, getJobByMonsterQuery, monID)
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
	SELECT j.id,j.started_at,j.updated_at,j.job_type, m.monster_id, j.user_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []getJobsQueryResult
	err := sqlscan.Select(ctx, c.conn.DB, &res, getJobByID, id)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}

	job := toJob(res)

	return &job, nil
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
	err := sqlscan.Select(ctx, c.conn.DB, &res, getJobsQuery)
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

func (c *Client) GetWoodcuttingJobByID(ctx context.Context, id int) (*WoodCuttingJob, error) {
	const getWoodcuttingJobByID = `
	SELECT j.id,j.user_id,j.started_at,j.updated_at, m.monster_id, w.tree_type
	FROM jobs as j
	LEFT JOIN woodcutting_jobs as w
	ON j.id=w.job_id
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []struct {
		job
		jobMonster
		TreeType string
	}
	err := sqlscan.Select(ctx, c.conn.DB, &res, getWoodcuttingJobByID, id)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	if len(res) > 1 {
		return nil, fmt.Errorf("multiple woodcutting jobs found")
	}
	woodCuttingJob := toWoodcuttingJob(res)
	return &woodCuttingJob, nil
}

func toJob(res []getJobsQueryResult) Job {
	job := Job{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, entry.MonsterID)
	}
	job.ID = res[0].ID
	job.UserID = res[0].UserID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	job.UpdatedAt = res[0].UpdatedAt
	return job
}

func toWoodcuttingJob(res []struct {
	job
	jobMonster
	TreeType string
},
) WoodCuttingJob {
	job := WoodCuttingJob{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, entry.MonsterID)
	}
	job.ID = res[0].ID
	job.UserID = res[0].UserID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	job.UpdatedAt = res[0].UpdatedAt
	job.TreeType = res[0].TreeType

	return job
}
