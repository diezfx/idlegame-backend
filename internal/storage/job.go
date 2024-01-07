package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/georgysavva/scany/sqlscan"
)

func (c *Client) StoreNewWoodCuttingJob(ctx context.Context, monsterID int, woodType string) (int, error) {
	var jobID int

	err := c.withTx(ctx, func(tx *sql.Tx) error {
		const insertJobQuery = `
		INSERT INTO jobs(started_at,job_type)
		values($1,$2)
		RETURNING id`

		err := tx.QueryRowContext(ctx, insertJobQuery, time.Now(), "woodcutting").Scan(&jobID)
		if err != nil {
			return err
		}

		const insertMonsterQuery = `
		INSERT INTO job_monsters (job_id, monster_id)
		values($1,$2)`
		_, err = tx.ExecContext(ctx, insertMonsterQuery, jobID, monsterID)
		if err != nil {
			return err
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
	// store all monsters
	// store specific details
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
	SELECT j.id,j.started_at,j.job_type, m.monster_id
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

func toJob(res []struct {
	job
	jobMonster
},
) Job {
	job := Job{}
	for _, entry := range res {
		job.Monsters = append(job.Monsters, MonsterID(entry.MonsterID))
	}
	job.ID = res[0].ID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	return job
}
