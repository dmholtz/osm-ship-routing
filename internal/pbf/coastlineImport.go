package pbf

import (
	"io"
	"log"
	"os"
	"runtime"

	cl "github.com/dmholtz/osm-ship-routing/pkg/coastline"
	"github.com/qedus/osmpbf"
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
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	//d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if o, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch o := o.(type) {
			case *osmpbf.Node:
				i.nodeCount++
				i.nodeIdMap[o.ID] = cl.NodeCoordinates{Lon: o.Lon, Lat: o.Lat}
			case *osmpbf.Way:
				if o.Tags["natural"] == "coastline" {
					coastline := cl.NewAtomicSegment(o.NodeIDs)
					i.coastlines = append(i.coastlines, coastline)
					i.coastlineCount++
				}
				i.wayCount++
			case *osmpbf.Relation:
			default:
				log.Fatalf("unknown type %T\n", o)
			}
		}
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
