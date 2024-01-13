package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
)

type Client struct {
	conn *postgres.DB
}

func New(ctx context.Context, sqlConn *postgres.DB) (*Client, error) {
	client := Client{conn: sqlConn}

	err := client.conn.Up(ctx)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}

	return &client, nil
}

func (c *Client) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
