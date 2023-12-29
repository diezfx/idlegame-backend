package storage

import (
	"context"
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
