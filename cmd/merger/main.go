package main

import (
	"fmt"
	"time"

	"github.com/dmholtz/osm-ship-routing/internal/pbf"
	"github.com/dmholtz/osm-ship-routing/pkg/coastline"
)

const pbfFile string = "antarctica.osm.pbf"
const geojsonFile string = "antarctica.geo.json"
const polyJsonFile string = "antarctica.poly.json"

//const pbfFile string = "central-america.osm.pbf"
//const geojsonFile string = "central-america.geo.json"

//const pbfFile string = "planet-coastlines.pbf"
//const geojsonFile string = "planet-coastlines.geo.json"

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
	fmt.Printf("Unmergable coastline segments: %d\n", merger.UnmergableSegmentCount())

	start = time.Now()

	//pbf.ExportGeojson(merger.Polygons(), coastlineImporter, geojsonFile)
	pbf.ExportPolygonJson(merger.Polygons(), coastlineImporter, polyJsonFile)

	elapsed = time.Since(start)
	fmt.Printf("[TIME] Export to geojson: %s\n", elapsed)
	fmt.Printf("Exported coastlines to %s\n", geojsonFile)
}
