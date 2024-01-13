package inventory

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

type InventoryStorage interface {
	AddItem(ctx context.Context, userID int, id string, quantity int) (int, error)
	GetItem(ctx context.Context, userID int, id string) (*storage.InventoryEntry, error)
	GetInventory(ctx context.Context, userID int) ([]storage.InventoryEntry, error)
}
