package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/georgysavva/scany/sqlscan"
)

func (s *Client) GetItem(ctx context.Context, userID int, itemID string) (*InventoryEntry, error) {
	const getItemQuantityQuery = `
	SELECT quantity,item_def_id,user_id
	FROM inventory_items
	WHERE user_id=$1 AND item_def_id=$2
	`
	var res InventoryEntry
	err := sqlscan.Get(ctx, s.conn, &res, getItemQuantityQuery, userID, itemID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get entry: %w", err)
	}
	return &res, nil
}

func (s *Client) AddItem(ctx context.Context, userID int, itemID string, quantity int) (int, error) {
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

	err := s.conn.QueryRowContext(ctx, addItemQuery, &newQuantity, addItemQuery, userID, itemID, quantity).Scan(&newQuantity)
	if err != nil {
		return 0, fmt.Errorf("add item: %w", err)
	}
	return newQuantity, nil
}

func (s *Client) GetInventory(ctx context.Context, userID int) ([]InventoryEntry, error) {
	var res []InventoryEntry
	const getInventoryQuery = `
	SELECT quantity,item_def_id,user_id
	FROM inventory_items
	WHERE user_id=$1
	`
	err := sqlscan.Select(ctx, s.conn.DB, &res, getInventoryQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("select inventory: %w", err)
	}
	return res, nil
}
