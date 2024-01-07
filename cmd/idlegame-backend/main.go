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

	srv, err := setup.SetupSplitService(ctx)
	if err != nil {
		logger.Fatal(context.Background(), err).Msg("failed setup")
	}

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(context.Background(), err).Msg("failed listening and serve")
	}
}
