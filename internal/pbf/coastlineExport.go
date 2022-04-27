package pbf

import (
	"bytes"
	"math"
	"os"

	cl "github.com/dmholtz/osm-ship-routing/pkg/coastline"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func roundSevenPlaces(number float64) float64 {
	return math.Floor(number*10e7) / 10e7
}

func ExportGeojson(coastlines []cl.AtomicSegment, importer *CoastlineImporter, filename string) {
	fc := geojson.NewFeatureCollection()

	for _, coastline := range coastlines {
		ls := orb.Ring(make([]orb.Point, 0, len(coastline)))
		for _, nodeId := range coastline {
			nc, err := importer.nodeIdMap[nodeId]
			if !err {
				panic(nodeId)
			}
			p := orb.Point{roundSevenPlaces(nc.Lon), roundSevenPlaces(nc.Lat)}
			ls = append(ls, p)
		}
		pol := orb.Polygon(make([]orb.Ring, 1, 1))
		pol[0] = ls
		f := geojson.NewFeature(pol)
		fc.Append(f)
	}

	jsonObj, _ := fc.MarshalJSON()
	jsonObj = bytes.ReplaceAll(jsonObj, []byte("null"), []byte("{}"))

	wErr := os.WriteFile(filename, jsonObj, 0644)
	if wErr != nil {
		panic(wErr)
	}
}
