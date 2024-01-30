package jobs

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/service/monster"
)

type MonsterService interface {
	GetMonsterByID(ctx context.Context, id int) (*monster.Monster, error)
	AddMonsterExperience(ctx context.Context, userID int, exp int) (int, error)
}
