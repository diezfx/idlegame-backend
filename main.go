package main

import (
	"context"
	"fmt"
	"log"

	"github.com/diezfx/idlegame-backend/internal/config"
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/postgres"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	psqlClient, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.New(ctx, psqlClient)
	if err != nil {
		log.Fatal(err)
	}

	res, err := store.GetJobByMonster(ctx, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", res)

	res2, err := store.GetJobByID(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", res2)
}
