package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/db"
	"github.com/golang-migrate/migrate/v4"
)

type Client struct {
	dbClient *db.DB
}

func New(ctx context.Context, sqlConn *db.DB) (*Client, error) {
	client := Client{dbClient: sqlConn}

	err := client.dbClient.Up(ctx)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("migrate up db: %w", err)
	}

	return &client, nil
}
