package jobs

import (
	"context"
	"time"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobStorage interface {
	StoreNewGatheringJob(ctx context.Context, jobType string, userID, monsterID int, woodType string) (int, error)
	DeleteGatheringJob(ctx context.Context, id int) error
	GetGatheringJobByID(ctx context.Context, id int) (*storage.GatheringJob, error)
	GetJobByMonster(ctx context.Context, monID int) (*storage.Job, error)
	GetJobByID(ctx context.Context, id int) (*storage.Job, error)
	UpdateJobUpdatedAt(ctx context.Context, id int, updatedAt time.Time) error
	GetJobs(ctx context.Context) ([]storage.Job, error)
}

type MonsterStorage interface {
	GetMonsterByID(ctx context.Context, id int) (*storage.Monster, error)
	AddMonsterExperience(ctx context.Context, userID int, exp int) (int, error)
}

type InventoryStorage interface {
	AddItems(ctx context.Context, items []storage.InventoryEntry) error
	GetItem(ctx context.Context, userID int, id string) (*storage.InventoryEntry, error)
	GetInventory(ctx context.Context, userID int) ([]storage.InventoryEntry, error)
}
