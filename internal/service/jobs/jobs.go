package jobs

import "github.com/diezfx/idlegame-backend/internal/service/item"

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

func New(jobStorage JobStorage, monsterStorage MonsterStorage, itemContainer *item.ItemContainer) *JobService {
	return &JobService{
		jobStorage: jobStorage, monsterStorage: monsterStorage,
		woodContainer: *InitWoodCutting(itemContainer),
	}
}

func (t JobType) String() string {
	return string(t)
}

type MonsterEntry struct {
	JobType   string
	MonsterID int
	JobID     int
}
