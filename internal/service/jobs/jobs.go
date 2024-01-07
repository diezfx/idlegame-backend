package jobs

// what kind of wood
type JobType string

const (
	WoodCuttingJobType JobType = "WoodCutting"
)

type JobService struct {
	jobStorage     JobStorage
	monsterStorage MonsterStorage
	woodContainer  WoodCuttingJobContainer
}

func (t JobType) String() string {
	return string(t)
}

type MonsterEntry struct {
	JobType   string
	MonsterID int
	JobID     int
}

func GetEntry(monID int) (MonsterEntry, error) {
}
