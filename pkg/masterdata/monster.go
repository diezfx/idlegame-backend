package masterdata

var monsterContainer MonsterContainer

type Monster struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Element MonsterElement `json:"element"`
}

type MonsterContainer struct {
	Monsters []Monster `json:"monsters"`
}

type MonsterElement string

func (t MonsterElement) String() string {
	return string(t)
}

const (
	FireType  MonsterElement = "Fire"
	WaterType MonsterElement = "Water"
	WindType  MonsterElement = "Wind"
	EarthType MonsterElement = "Earth"
)

func GetMonster(id int) Monster {
	return monsterContainer.Monsters[id]
}
