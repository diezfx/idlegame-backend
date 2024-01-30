package jobs

import (
	"context"
	"time"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

type JobStorage interface {
	StoreNewJob(ctx context.Context, jobType string, userID, monsterID int, jobDefID string) (int, error)
	GetJobByMonster(ctx context.Context, monID int) (*storage.Job, error)
	GetJobByID(ctx context.Context, id int) (*storage.Job, error)
	DeleteJobByID(ctx context.Context, id int) error
	UpdateJobUpdatedAt(ctx context.Context, id int, updatedAt time.Time) error
	GetJobs(ctx context.Context) ([]storage.Job, error)
}

type InventoryStorage interface {
	AddItems(ctx context.Context, items []storage.InventoryEntry) error
	GetItem(ctx context.Context, userID int, id string) (*storage.InventoryEntry, error)
	GetInventory(ctx context.Context, userID int) ([]storage.InventoryEntry, error)
}
