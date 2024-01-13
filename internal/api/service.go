package api

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/jobs"
)

type JobService interface {
	StartWoodCuttingJob(ctx context.Context, userID, monsterID int, treeType jobs.TreeType) (int, error)
	GetJob(ctx context.Context, id int) (*jobs.Job, error)
	GetWoodcuttingJob(ctx context.Context, id int) (*jobs.WoodCuttingJob, error)
	StopWoodCuttingJob(ctx context.Context, id int) error
}

type InventoryService interface {
	AddItem(ctx context.Context, userID int, id string, quantity int) (int, error)
	RemoveItem(ctx context.Context, userID int, id string, quantity int) (int, error)
	GetInventory(ctx context.Context, userID int) (*inventory.Inventory, error)
	GetItem(ctx context.Context, userID int, id string) (*inventory.Item, error)
}
