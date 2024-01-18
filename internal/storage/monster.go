package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/db"
	"github.com/diezfx/idlegame-backend/pkg/logger"
)

func (c *Client) GetMonsterByID(ctx context.Context, id int) (*Monster, error) {
	getMonsterQuery := `
	SELECT id,name,experience,monster_def_id
	FROM monsters
	WHERE id=$1
	`
	var monster Monster

	err := c.dbClient.Get(ctx, &monster, getMonsterQuery, id)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get monster: %w", err)
	}
	return &monster, nil
}

func (c *Client) AddMonster(ctx context.Context, mon Monster) (*Monster, error) {
	const addMonsterQuery = `
	INSERT INTO monsters(name,experience,monster_def_id)
	VALUES($1,$2,$3)
	RETURNING id
	`
	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		err := tx.Get(ctx, &mon.ID, addMonsterQuery, mon.Name, mon.Experience, mon.MonsterDefID)
		if err != nil {
			return fmt.Errorf("insert monster: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("execute transaction: %w", err)
	}
	return &mon, nil
}

func (c *Client) AddMonsterExperience(ctx context.Context, monID, additionalExp int) (int, error) {
	var currentExp int
	const addMonsterExperienceQuery = `
	UPDATE monsters
	SET experience=experience+$1
	WHERE id=$2
	RETURNING experience
	`
	err := c.dbClient.WithTx(ctx, func(tx db.Querier) error {
		err := tx.Get(ctx, &currentExp, addMonsterExperienceQuery, additionalExp, monID)
		if err != nil {
			return fmt.Errorf("add experience: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("execute transaction: %w", err)
	}
	logger.Debug(ctx).Any("exp", currentExp).Msg("added experience")
	return currentExp, nil
}
