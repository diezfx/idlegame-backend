package monster

import (
	"context"

	"github.com/diezfx/idlegame-backend/internal/storage"
)

type MonsterStorage interface {
	GetMonsterByID(ctx context.Context, id int) (*storage.Monster, error)
	AddMonsterExperience(ctx context.Context, userID int, exp int) (int, error)
}
