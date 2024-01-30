package monster

import (
	"context"
	"fmt"

	"github.com/diezfx/idlegame-backend/pkg/logger"
)

type Service struct {
	monsterStorage MonsterStorage
}

func New(monsterStorage MonsterStorage) *Service {
	return &Service{
		monsterStorage: monsterStorage,
	}
}

func (s *Service) GetMonsterByID(ctx context.Context, id int) (*Monster, error) {
	storageMon, err := s.monsterStorage.GetMonsterByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get monster %d: %w", id, err)
	}
	return MonsterFromStorage(storageMon), nil
}

func (s *Service) AddMonsterExperience(ctx context.Context, userID int, exp int) (int, error) {
	logger.Debug(ctx).Int("exp", exp).Msg("add experience")
	return s.monsterStorage.AddMonsterExperience(ctx, userID, exp)
}
