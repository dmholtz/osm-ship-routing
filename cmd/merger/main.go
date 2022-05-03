package main

import (
	"fmt"
	"time"

	"github.com/dmholtz/osm-ship-routing/internal/pbf"
	"github.com/dmholtz/osm-ship-routing/pkg/coastline"
	//"github.com/dmholtz/osm-ship-routing/pkg/geometry"
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
	/*
		p1 := Point2D{80, 80}
		p2 := Point2D{80, 40}
		p3 := Point2D{100, 60}
		p4 := Point2D{100, 10}
		p5 := Point2D{10, 10}
		p6 := Point2D{10, 80}
		p7 := Point2D{80, 80}
		testPointsInside := []Point2D{{70, 40}}
		testPointsOutside := []Point2D{{0, 50}}

		points := []*Point2D{&p1, &p2, &p3, &p4, &p5, &p6, &p7}
		polygon := Polygon2D{points}

		for _, testPoint := range testPointsInside {
			isInPolygon := polygon.IsInPolygon(&testPoint)
			if !isInPolygon {
				panic(fmt.Sprintf("Point is declared wrongly: %v", testPoint))
			}
		}
		for _, testPoint := range testPointsOutside {
			isInPolygon := polygon.IsInPolygon(&testPoint)
			if isInPolygon {
				panic(fmt.Sprintf("Point is declared wrongly: %v", testPoint))
			}
		}
	*/
}
