package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diezfx/idlegame-backend/pkg/postgres"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/golang-migrate/migrate/v4"
)

type Client struct {
	conn *postgres.DB
}

func New(ctx context.Context, sqlConn *postgres.DB) (*Client, error) {
	client := Client{conn: sqlConn}

	err := client.conn.Up(ctx)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}

	return &client, nil
}

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

	job := ToJob(res)

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

	job := ToJob(res)

	return &job, nil
}

func ToJob(res []struct {
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

func (c *Client) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
