package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/db"
)

func (c *Client) GetItem(ctx context.Context, userID int, itemID string) (*InventoryEntry, error) {
	const getItemQuantityQuery = `
	SELECT quantity,item_def_id,user_id
	FROM inventory_items
	WHERE user_id=$1 AND item_def_id=$2
	`
	var res InventoryEntry
	err := c.dbClient.Get(ctx, &res, getItemQuantityQuery, userID, itemID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get entry: %w", err)
	}
	return &res, nil
}

func (c *Client) AddItem(ctx context.Context, userID int, itemID string, quantity int) (int, error) {
	const addItemQuery = `
	INSERT INTO inventory_items
	(user_id, item_def_id, quantity)
	VALUES
	($1, $2, $3)
	ON CONFLICT (user_id, item_def_id) DO UPDATE
	SET quantity = inventory.quantity + $3
	RETURNING quantity
	`
	var newQuantity int

	err := c.dbClient.Get(ctx, &newQuantity, addItemQuery, addItemQuery, userID, itemID, quantity)
	if err != nil {
		return 0, fmt.Errorf("add item: %w", err)
	}
	return newQuantity, nil
}

func (c *Client) AddItems(ctx context.Context, items []InventoryEntry) error {
	const addItemQuery = `
	INSERT INTO inventory_items
	(user_id, item_def_id, quantity)
	VALUES
	($1, $2, $3)
	ON CONFLICT (user_id, item_def_id) DO UPDATE
	SET quantity = inventory_items.quantity + $3
	RETURNING quantity
	`

	return c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		var newQuantity int
		for _, item := range items {
			err := tx.Get(ctx, &newQuantity, addItemQuery, item.UserID, item.ItemDefID, item.Quantity)
			if err != nil {
				return fmt.Errorf("add item %s: %w", item.ItemDefID, err)
			}
		}
		return nil
	})
}

func (c *Client) GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error) {
	var res []InventoryEntry
	const getInventoryQuery = `
	SELECT quantity,item_def_id,user_id
	FROM inventory_items
	WHERE user_id=$1
	`
	err := c.dbClient.Select(ctx, &res, getInventoryQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("select inventory: %w", err)
	}
	return res, nil
}
