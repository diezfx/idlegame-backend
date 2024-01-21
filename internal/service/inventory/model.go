package inventory

import "github.com/diezfx/idlegame-backend/internal/storage"

type Inventory struct {
	UserID int
	Items  []Item
}

type Item struct {
	Quantity  int
	ItemDefID string
}

func ToItemFromStorage(i storage.InventoryEntry) *Item {
	return &Item{
		Quantity:  i.Quantity,
		ItemDefID: i.ItemDefID,
	}
}

func ToInventoryFromStorageEntries(entries []storage.InventoryEntry, userID int) *Inventory {
	items := make([]Item, 0, len(entries))
	for _, entry := range entries {
		items = append(items, *ToItemFromStorage(entry))
	}
	return &Inventory{
		UserID: userID,
		Items:  items,
	}
}
