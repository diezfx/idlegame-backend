package storage

import (
	"context"
	"fmt"
	"time"
)

const insertJobQuery = `
INSERT INTO jobs(user_id,started_at,updated_at,job_type)
values($1,$2,$2,$3)
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
	SELECT j.id,j.started_at,j.updated_at,j.job_type, m.monster_id, j.user_id
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
	job.UserID = res[0].UserID
	job.JobType = res[0].JobType
	job.StartedAt = res[0].StartedAt
	job.UpdatedAt = res[0].UpdatedAt
	return job
}
