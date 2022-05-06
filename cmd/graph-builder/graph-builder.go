package main

import (
	"encoding/json"
	"os"

	grid "github.com/dmholtz/osm-ship-routing/pkg/grid"
)

func main() {

	sgg := grid.NewSphereGridGraph(40, 20)

	jsonObj, err := json.Marshal(sgg)
	if err != nil {
		panic(err)
	}

	wErr := os.WriteFile("graph.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}
}
