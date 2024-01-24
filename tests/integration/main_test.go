//go:build !unit

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/pkg/db"
	"github.com/diezfx/idlegame-backend/pkg/logger"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	configLoader, err := config.Load()
	if err != nil {
		panic(err)
	}

	dbClient, err := db.New(configLoader.DB)
	if err != nil {
		panic(err)
	}

	cleanupFunc := seedDb(ctx, dbClient)
	defer func() {
		cleanupFunc()
		if err := recover(); err != nil {
			fmt.Println(err)
		}

	}()
	// setup
	m.Run()
	// teardown
}

func seedDb(ctx context.Context, db *db.DB) func() {

	seedFile, err := os.ReadFile("db/seed.sql")
	if err != nil {
		logger.Fatal(ctx, err).Msg("read seed file")
	}

	cleanFile, err := os.ReadFile("db/clean.sql")
	if err != nil {
		logger.Fatal(ctx, err).Msg("read clean file")
	}

	_, err = db.Exec(ctx, string(cleanFile))
	if err != nil {
		logger.Fatal(ctx, err).Msg("exec cleaning")
	}
	_, err = db.Exec(ctx, string(seedFile))
	if err != nil {
		logger.Fatal(ctx, err).Msg("exec seeding")
	}

	return func() {
		_, err = db.Exec(ctx, string(cleanFile))
		if err != nil {
			logger.Fatal(ctx, err).Msg("exec cleaning")
		}
	}
}
