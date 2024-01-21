package api

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/inventory"
	"github.com/diezfx/idlegame-backend/internal/service/item"
	"github.com/diezfx/idlegame-backend/internal/service/jobs"
	"github.com/diezfx/idlegame-backend/internal/service/monster"
)

type JobService interface {
	GetJob(ctx context.Context, id int) (*jobs.Job, error)
	StopJob(ctx context.Context, id int) error

	StartWoodCuttingJob(ctx context.Context, userID, monsterID int, treeType item.TreeType) (int, error)
	GetWoodcuttingJob(ctx context.Context, id int) (*jobs.WoodCuttingJob, error)
	StartMiningJob(ctx context.Context, userID, monsterID int, oreType item.OreType) (int, error)
	StartHarvestingJob(ctx context.Context, userID, monsterID int, cropType item.CropType) (int, error)
	GetHarvestingJob(ctx context.Context, id int) (*jobs.HarvestingJob, error)
	GetMiningJob(ctx context.Context, id int) (*jobs.MiningJob, error)
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
