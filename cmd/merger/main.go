package main

import (
	"fmt"
	"time"

	"github.com/dmholtz/osm-ship-routing/internal/pbf"
	"github.com/dmholtz/osm-ship-routing/pkg/coastline"
)

//const pbfFile string = "antarctica.osm.pbf"

const pbfFile string = "planet-coastlines.pbf"

func main() {

	start := time.Now()

	coastlineImporter := pbf.NewImporter(pbfFile)
	coastlineImporter.Import()

	elapsed := time.Since(start)
	fmt.Printf("[TIME] Import: %s\n", elapsed)

	start = time.Now()

	merger := coastline.NewMerger(coastlineImporter.Coastlines())
	merger.Merge()

	elapsed = time.Since(start)
	fmt.Printf("[TIME] Merge: %s\n", elapsed)
	fmt.Printf("Polygon Count: %d\n", len(merger.Polygons()))
	fmt.Printf("Merge Count: %d\n", merger.MergeCount())
}
