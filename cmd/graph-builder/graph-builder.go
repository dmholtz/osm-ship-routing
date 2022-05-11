package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	"github.com/dmholtz/osm-ship-routing/pkg/graph"
	"github.com/dmholtz/osm-ship-routing/pkg/grid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

const density = 20

func main() {

	//arg := africPolygons()
	arg := loadGeoJsonPolygons("antarctica.geo.json")
	//arg := loadGeoJsonPolygons("planet-coastlines.geo.json")

	sgg := grid.NewSimpleSphereGrid(2*density, density, arg)

	gridGraph := sgg.ToGraph()
	jsonObj, err := json.Marshal(gridGraph)
	if err != nil {
		panic(err)
	}

	wErr := os.WriteFile("graph.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}

	aag := graph.NewAdjacencyArrayFromGraph(gridGraph)
	jsonObj, err = json.Marshal(aag)
	if err != nil {
		panic(err)
	}

	wErr = os.WriteFile("adjacency_array_graph.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}
}

func loadGeoJsonPolygons(file string) []geometry.Polygon {

	start := time.Now()
	json, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("[TIME] Read file: %s\n", elapsed)

	polygons := make([]geometry.Polygon, 0)

	start = time.Now()
	fc, _ := geojson.UnmarshalFeatureCollection(json)
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Unmarshal: %s\n", elapsed)

	start = time.Now()
	for _, f := range fc.Features {
		pol := f.Geometry.(orb.Polygon)
		pts := pol[0]

		points := make([]*geometry.Point, 0)
		for _, pt := range pts {
			points = append(points, geometry.NewPoint(pt[1], pt[0]))
		}
		polygon := geometry.NewPolygon(points)
		polygons = append(polygons, *polygon)
	}
	elapsed = time.Since(start)
	fmt.Printf("[TIME] Convert to internal polygon type: %s\n", elapsed)

	return polygons
}

func africPolygons() []geometry.Polygon {
	var africa [][][]float64 = [][][]float64{
		{
			{
				-4.5703125,
				35.17380831799959,
			},
			{
				-9.84375,
				32.84267363195431,
			},
			{
				-17.578125,
				20.96143961409684,
			},
			{
				-16.875,
				12.554563528593656,
			},
			{
				-8.0859375,
				4.214943141390651,
			},
			{
				3.8671874999999996,
				6.664607562172573,
			},
			{
				9.84375,
				4.214943141390651,
			},
			{
				8.4375,
				0,
			},
			{
				13.0078125,
				-5.615985819155327,
			},
			{
				13.0078125,
				-19.973348786110602,
			},
			{
				19.6875,
				-34.30714385628803,
			},
			{
				29.179687499999996,
				-33.43144133557529,
			},
			{
				30.937499999999996,
				-25.799891182088306,
			},
			{
				35.859375,
				-19.642587534013032,
			},
			{
				40.078125,
				-15.28418511407642,
			},
			{
				40.078125,
				-4.915832801313164,
			},
			{
				51.67968749999999,
				10.487811882056695,
			},
			{
				42.5390625,
				8.407168163601076,
			},
			{
				31.289062500000004,
				29.84064389983441,
			},
			{
				21.796875,
				32.54681317351514,
			},
			{
				18.6328125,
				30.44867367928756,
			},
			{
				10.546875,
				33.137551192346145,
			},
			{
				8.0859375,
				37.71859032558816,
			},
			{
				-4.5703125,
				35.17380831799959,
			},
		},
	}

	arg := make([]geometry.Polygon, 0)
	for _, pol := range africa {
		myPoints := make([]*geometry.Point, 0)
		for _, point := range pol {
			myPoint := geometry.NewPoint(point[1], point[0])
			myPoints = append(myPoints, myPoint)
		}
		myPol := geometry.NewPolygon(myPoints)
		arg = append(arg, *myPol)
	}
	return arg
}
