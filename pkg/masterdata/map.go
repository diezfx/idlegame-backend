package masterdata

import (
	"fmt"

	"github.com/lafriks/go-tiled"
)

type Map struct {
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	TileSet []Tile  `json:"tileSet"`
	Layer   [][]int `json:"layer"`
}

type Tile struct {
	TileSetID int    `json:"tileSetId"`
	Source    string `json:"source"`
	Type      string `json:"type"`
}

func LoadMap(cfg Config) (*Map, error) {

	gameMap, err := tiled.LoadFile(cfg.MapPath)
	if err != nil {
		return nil, fmt.Errorf("load map from file: %w", err)
	}

	tileSet := make([]Tile, 0, len(gameMap.Tilesets[0].Tiles))
	for _, tile := range gameMap.Tilesets[0].Tiles {
		tileSet = append(tileSet, Tile{
			TileSetID: int(tile.ID),
			Source:    tile.Image.Source,
			Type:      tile.Properties.GetString("Type"),
		})
	}

	layer := make([][]int, gameMap.Height)
	for y := 0; y < gameMap.Width; y++ {
		layer[y] = make([]int, gameMap.Width)
		for x := 0; x < gameMap.Height; x++ {
			layer[y][x] = int(gameMap.Layers[0].Tiles[y*gameMap.Height+x].ID)
		}
	}

	return &Map{
		Width:   gameMap.Width,
		Height:  gameMap.Height,
		TileSet: tileSet,
		Layer:   layer,
	}, nil

}
