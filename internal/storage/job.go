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
		INSERT INTO jobs(user_id,started_at,job_type)
		values($1,$2,$3)
		RETURNING id`

		err := tx.QueryRowContext(ctx, insertJobQuery, userID, time.Now(), "woodcutting").Scan(&jobID)
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

func (c *Client) GetJobByMonster(ctx context.Context, monID int) (*Job, error) {
	const getJobByMonsterQuery = `
	SELECT j.id,j.started_at,j.job_type, m.monster_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE m.monster_id=$1`

	var res []struct {
		job
		jobMonster
	}
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
	SELECT j.id,j.started_at,j.job_type, m.monster_id, j.user_id
	FROM jobs as j
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []struct {
		job
		jobMonster
	}
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

func (c *Client) GetWoodcuttingJobByID(ctx context.Context, id int) (*WoodCuttingJob, error) {
	const getWoodcuttingJobByID = `
	SELECT j.id,j.started_at, m.monster_id, w.tree_type
	FROM jobs as j
	LEFT JOIN woodcutting_jobs as w
	ON j.id=w.job_id
	LEFT JOIN job_monsters as m
	ON j.id=m.job_id
	WHERE j.id=$1`

	var res []WoodCuttingJob
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
	return &res[0], nil
}

func toJob(res []struct {
	job
	jobMonster
},
) Job {
	job := Job{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, entry.MonsterID)
	}
	job.ID = res[0].ID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	return job
}
