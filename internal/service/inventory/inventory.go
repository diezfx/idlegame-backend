package inventory

import (
	"context"
	"fmt"

	"github.com/diezfx/idlegame-backend/internal/service"
)

type Service struct {
	storage InventoryStorage
}

func New(storage InventoryStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddItem(ctx context.Context, userID int, itemID string, quantity int) (int, error) {
	return s.storage.AddItem(ctx, userID, itemID, quantity)
}

func (s *Service) RemoveItem(ctx context.Context, userID int, itemID string, removeQuantity int) (int, error) {
	currentQuantity, err := s.storage.GetItem(ctx, userID, itemID)
	if err != nil {
		return 0, fmt.Errorf("get quantity: %w", err)
	}
	if currentQuantity.Quantity < removeQuantity {
		return 0, service.ErrNotEnoughInventory
	}

	return s.storage.AddItem(ctx, userID, itemID, -removeQuantity)
}

func (s *Service) GetItem(ctx context.Context, userID int, itemID string) (*Item, error) {
	item, err := s.storage.GetItem(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	return ToItemFromStorage(*item), nil
}

func (s *Service) GetInventory(ctx context.Context, userID int) (*Inventory, error) {
	res, err := s.storage.GetInventory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get inventory: %w", err)
	}
	return ToInventoryFromStorageEntries(res, userID), nil
}
