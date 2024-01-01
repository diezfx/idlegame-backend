package jobs

//what kind of wood

type TreeType string

const (
	SpruceType TreeType = "Spruce"
	BirchType  TreeType = "Birch"
	PineType   TreeType = "Pine"
)

type WoodCuttingJob struct {
	ID       int
	Monster  int
	TreeType TreeType
}

// different tree types have
// level requirements, durations, exp gains

type WoodCuttingJobDefinition struct {
	ExpGain          int
	Reward           TreeType
	LevelRequirement int
}
