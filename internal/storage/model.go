package storage

import "time"

type MonsterID int

type jobMonster struct {
	MonsterID int
	JobID     int
}

type job struct {
	ID        int
	StartedAt time.Time
	JobType   string
}

type WoodCuttingJob struct {
	ID        int
	Monster   int
	TreeType  string
	StartedAt time.Time
}

type Job struct {
	ID        int
	StartedAt time.Time
	Monsters  []MonsterID
	JobType   string
}

type Monster struct {
	ID           int
	Name         string
	Experience   int
	MonsterDefID int
}
