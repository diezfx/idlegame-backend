package monster

type MonsterType string

const (
	FireType  MonsterType = "Fire"
	WaterType MonsterType = "Water"
	WindType  MonsterType = "Wind"
	EarthType MonsterType = "Earth"
)

var monList = map[int]MonsterDefinition{
	1: {ID: 1, Name: "schiggo", Type: WaterType},
	2: {ID: 2, Name: "bisa", Type: EarthType},
	3: {ID: 3, Name: "glumander", Type: FireType}}

type MonsterDefinition struct {
	ID   int
	Name string
	Type MonsterType
}

type Monster struct {
	ID           int
	MonsterDefID int
	Experience   int
}

func New(id, monsterDefID int) Monster {
	return Monster{
		ID:           id,
		MonsterDefID: monsterDefID,
		Experience:   0,
	}
}

func (m *Monster) AddExperience(exp int) {
	m.Experience += exp
}

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

func (m *Monster) Type() MonsterType {
	return monList[m.MonsterDefID].Type
}
