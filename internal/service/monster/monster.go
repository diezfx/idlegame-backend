package monster

import (
	"github.com/diezfx/idlegame-backend/internal/storage"
	"github.com/diezfx/idlegame-backend/pkg/masterdata"
)

type Monster struct {
	ID           int
	MonsterDefID int
	Experience   int
}

func NewMonster(id, monsterDefID int) Monster {
	return Monster{
		ID:           id,
		MonsterDefID: monsterDefID,
		Experience:   0,
	}
}

func (m *Monster) AddExperience(exp int) {
	m.Experience += exp
}

//nolint:gomnd // fine for level
func (m *Monster) Level() int {
	if m.Experience < 100 {
		return 1
	}
	if m.Experience < 200 {
		return 2
	}
	if m.Experience < 400 {
		return 3
	}
	if m.Experience < 800 {
		return 4
	}
	return 5
}

func (m *Monster) Element() masterdata.MonsterElement {
	return masterdata.GetMonster(m.MonsterDefID).Element
}

func (m *Monster) Name() string {
	return masterdata.GetMonster(m.MonsterDefID).Name
}

func MonsterFromStorage(m *storage.Monster) *Monster {
	return &Monster{
		ID:           m.ID,
		MonsterDefID: m.MonsterDefID,
		Experience:   m.Experience,
	}
}
