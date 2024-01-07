package jobs

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/monster"
)

type JobStorage interface {
	StoreMonsterEntry(ctx context.Context, entry MonsterEntry) error
	GetMonsterEntry(ctx context.Context, monID int) (*MonsterEntry, error)

	StoreWoodCuttingJob(ctx context.Context, job WoodCuttingJob) error
	GetWoodCuttingJob(ctx context.Context, id int) (*WoodCuttingJob, error)
}

type MonsterStorage interface {
	GetMonster(id int) (*monster.Monster, error)
}
