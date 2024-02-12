package main

import (
	"fmt"
	"os"

	"github.com/lafriks/go-tiled"
)

func main() {
	const mapPath = "assets/map.tmx"
	gameMap, err := tiled.LoadFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	for _, tile := range gameMap.Tilesets[0].Tiles {
		fmt.Println(tile.Properties.GetString("Type"))
	}

	fmt.Println(gameMap)
}
