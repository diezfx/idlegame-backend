package storage

import "time"

type jobMonster struct {
	MonsterID int
	JobID     int
}

type job struct {
	ID        int
	JobDefID  string
	StartedAt time.Time
	UpdatedAt time.Time
	JobType   string
	UserID    int
}

type Job struct {
	ID        int
	JobDefID  string
	UserID    int
	StartedAt time.Time
	UpdatedAt time.Time
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
