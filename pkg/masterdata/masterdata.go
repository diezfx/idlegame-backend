package masterdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Container struct {
	Monsters MonsterContainer
	Jobs     JobContainer
	Items    ItemContainer
}

type Config struct {
	Path string
}

func New(cfg Config) (*Container, error) {
	// load items
	itemPath := path.Join(cfg.Path, "items.json")
	itemFile, err := os.ReadFile(itemPath)
	if err != nil {
		return nil, fmt.Errorf("read item file %s: %w", itemPath, err)
	}
	var items ItemContainer
	err = json.Unmarshal(itemFile, &items)
	if err != nil {
		return nil, fmt.Errorf("unmarshal items: %w", err)
	}
	// load jobs

	jobPath := path.Join(cfg.Path, "jobs.json")
	jobFile, err := os.ReadFile(jobPath)
	if err != nil {
		return nil, fmt.Errorf("read job file %s: %w", jobPath, err)
	}

	var jobs JobContainer
	err = json.Unmarshal(jobFile, &jobs)
	if err != nil {
		return nil, fmt.Errorf("unmarshal jobs: %w", err)
	}

	// load monsters
	monsterPath := path.Join(cfg.Path, "monsters.json")
	monsterFile, err := os.ReadFile(monsterPath)
	if err != nil {
		return nil, fmt.Errorf("read monster file %s: %w", monsterPath, err)
	}
	var monsters MonsterContainer
	err = json.Unmarshal(monsterFile, &monsters)
	if err != nil {
		return nil, fmt.Errorf("unmarshal monsters: %w", err)
	}
	monsterContainer = monsters

	return &Container{
		Monsters: monsters,
		Jobs:     jobs,
		Items:    items,
	}, nil
}
