package pbf

import (
	"context"
	"log"
	"os"
	"runtime"

	cl "github.com/dmholtz/osm-ship-routing/pkg/coastline"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

type CoastlineImporter struct {
	pbfFile        string
	nodeIdMap      map[int64]cl.NodeCoordinates
	coastlines     []cl.Segment
	nodeCount      int64
	wayCount       int64
	coastlineCount int64
}

func NewImporter(pbfFile string) *CoastlineImporter {
	nodeIdMap := make(map[int64]cl.NodeCoordinates, 80638687) // pre-allocation improves performance roughly by 1.8
	coastlines := []cl.Segment{}
	return &CoastlineImporter{pbfFile: pbfFile, nodeIdMap: nodeIdMap, coastlines: coastlines}
}

func (i *CoastlineImporter) Import() {
	f, err := os.Open(i.pbfFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, runtime.GOMAXPROCS(-1))
	defer scanner.Close()

	scanner.SkipRelations = true

	for scanner.Scan() {
		switch o := scanner.Object().(type) {

		case *osm.Node:
			i.nodeCount++
			i.nodeIdMap[int64(o.ID)] = cl.NodeCoordinates{Lon: o.Lon, Lat: o.Lat}
		case *osm.Way:
			for _, tag := range o.Tags {
				if tag.Key == "natural" && tag.Value == "coastline" {
					raw := make([]int64, len(o.Nodes), len(o.Nodes))
					for i, node := range o.Nodes {
						raw[i] = int64(node.ID)
					}
					coastline := cl.NewAtomicSegment([]int64(raw))
					i.coastlines = append(i.coastlines, coastline)
					i.coastlineCount++
					break
				}
			}
			i.wayCount++
		case *osm.Relation:
		default:
			log.Fatalf("unknown type %T\n", o)
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}
}

func (i *CoastlineImporter) Coastlines() []cl.Segment {
	return i.coastlines
}

func (i *CoastlineImporter) String() string {
	/*
		fmt.Printf("Node count: %d\n", nodeCount)
		fmt.Printf("NodeMap size: %d\n", len(nodes))
		fmt.Printf("Way count: %d\n", waysCount)
		fmt.Printf("Coastline count: %d\n", coastlineCount)
		fmt.Printf("Coastline check: %d\n", len(coastlines))
	*/
	return ""
}
