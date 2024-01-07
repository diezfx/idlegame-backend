package jobs

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/monster"
	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobStorage interface {
	StoreNewWoodCuttingJob(ctx context.Context, monsterID int, woodType string) (int, error)
	GetJobByMonster(ctx context.Context, monID int) (*storage.Job, error)
	GetJobByID(ctx context.Context, id int) (*storage.Job, error)
}

type MonsterStorage interface {
	GetMonster(id int) (*monster.Monster, error)
}
