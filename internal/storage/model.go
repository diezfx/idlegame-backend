package storage

import "time"

type jobMonster struct {
	MonsterID int
	JobID     int
}

type job struct {
	ID        int
	StartedAt time.Time
	JobType   string
	UserID    int
}

type WoodCuttingJob struct {
	ID        int
	MonsterID int
	TreeType  string
	StartedAt time.Time
}

type Job struct {
	ID        int
	StartedAt time.Time
	Monsters  []int
	JobType   string
}

type Monster struct {
	ID           int
	Name         string
	Experience   int
	MonsterDefID int
}

type InventoryEntry struct {
	Quantity  int
	UserID    int
	ItemDefID string
}
