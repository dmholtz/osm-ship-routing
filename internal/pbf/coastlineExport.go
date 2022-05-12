package pbf

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"strconv"

	cl "github.com/dmholtz/osm-ship-routing/pkg/coastline"
	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func roundSevenPlaces(number float64) float64 {
	return math.Floor(number*10e7) / 10e7
}

func ExportGeojson(coastlines []cl.AtomicSegment, importer *CoastlineImporter, filename string) {
	fc := geojson.NewFeatureCollection()

	for _, coastline := range coastlines {
		ring := orb.Ring(make([]orb.Point, 0, len(coastline)))
		for _, nodeId := range coastline {
			nc, in := importer.nodeIdMap[nodeId]
			if !in {
				panic("NodeId " + strconv.FormatInt(nodeId, 10) + " not in map of imported nodes\n")
			}
			p := orb.Point{roundSevenPlaces(nc.Lon), roundSevenPlaces(nc.Lat)}
			ring = append(ring, p)
		}
		pol := orb.Polygon(make([]orb.Ring, 1, 1))
		pol[0] = ring
		f := geojson.NewFeature(pol)
		fc.Append(f)
	}

	jsonObj, _ := fc.MarshalJSON()
	// replace json 'null' value with '{}' to ensure compatibility with geojson.io
	jsonObj = bytes.ReplaceAll(jsonObj, []byte("null"), []byte("{}"))

	wErr := os.WriteFile(filename, jsonObj, 0644)
	if wErr != nil {
		panic(wErr)
	}
}

func ExportPolygonJson(coastlines []cl.AtomicSegment, importer *CoastlineImporter, filename string) {
	polygons := make([]geometry.Polygon, 0)

	for _, coastline := range coastlines {
		polygon := make(geometry.Polygon, 0)
		for _, nodeId := range coastline {
			nc, ok := importer.nodeIdMap[nodeId]
			if !ok {
				panic("NodeId " + strconv.FormatInt(nodeId, 10) + " not in map of imported nodes\n")
			}
			p := geometry.NewPoint(roundSevenPlaces(nc.Lat), roundSevenPlaces(nc.Lon))
			polygon = append(polygon, p)
		}
		polygons = append(polygons, polygon)
	}

	jsonObj, _ := json.Marshal(polygons)

	wErr := os.WriteFile(filename, jsonObj, 0644)
	if wErr != nil {
		panic(wErr)
	}
}
