package main

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/setup"
	"github.com/diezfx/idlegame-backend/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		logger.Debug(ctx).Err(err).Msg("Error loading .env file")
	}

	srv, daemon, err := setup.SetupSplitService(ctx)
	if err != nil {
		logger.Fatal(ctx, err).Msg("failed setup")
	}

	go func() {
		err := daemon.Run(ctx)
		if err != nil {
			logger.Fatal(ctx, err).Msg("failed daemon")
		}
	}()

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(context.Background(), err).Msg("failed listening and serve")
	}
}
