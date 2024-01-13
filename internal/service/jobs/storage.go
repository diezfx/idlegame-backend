package jobs

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobStorage interface {
	StoreNewWoodCuttingJob(ctx context.Context, userID, monsterID int, woodType string) (int, error)
	DeleteWoodCuttingJob(ctx context.Context, id int) error
	GetWoodcuttingJobByID(ctx context.Context, id int) (*storage.WoodCuttingJob, error)
	GetJobByMonster(ctx context.Context, monID int) (*storage.Job, error)
	GetJobByID(ctx context.Context, id int) (*storage.Job, error)
}

type MonsterStorage interface {
	GetMonsterByID(ctx context.Context, id int) (*storage.Monster, error)
}
