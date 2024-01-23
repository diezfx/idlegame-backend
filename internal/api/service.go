package api

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/jobs"
	"github.com/diezfx/idlegame-backend/internal/service/monster"
)

type JobService interface {
	GetJob(ctx context.Context, id int) (*jobs.Job, error)
	GetJobs(ctx context.Context) ([]jobs.Job, error)
	StopJob(ctx context.Context, id int) error

	StartGatheringJob(ctx context.Context, userID, monsterID int, jobDefID string) (int, error)
	StartProcessingJob(ctx context.Context, userID, monsterID int, jobDefID string) (int, error)
}

type InventoryService interface {
	AddItem(ctx context.Context, userID int, id string, quantity int) (int, error)
	RemoveItem(ctx context.Context, userID int, id string, quantity int) (int, error)
	GetInventory(ctx context.Context, userID int) (*inventory.Inventory, error)
	GetItem(ctx context.Context, userID int, id string) (*inventory.Item, error)
}

type MonsterService interface {
	GetMonsterByID(ctx context.Context, id int) (*monster.Monster, error)
}
