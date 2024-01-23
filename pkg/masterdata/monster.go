package masterdata

type Monster struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type MonsterContainer struct {
	Monsters []Monster `json:"monsters"`
}
