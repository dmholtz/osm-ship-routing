package main

import (
	"encoding/json"
	"os"

	"github.com/dmholtz/osm-ship-routing/pkg/geometry"
	"github.com/dmholtz/osm-ship-routing/pkg/grid"
)

const density = 80

func main() {

	var polygons1 [][][]float64 = [][][]float64{
		{
			{-78.34941069014627, -30.234375},
			{-77.76758238272801, -57.65624999999999},
			{-75.67219739055291, -126.91406249999999},
			{-81.03861703916249, -163.4765625},
			{-80.05804956215623, 160.3125},
			{-69.162557908105, 149.0625},
			{-71.41317683396565, 11.6015625},
			{-78.34941069014627, -30.234375},
		},
	}
	polygons1 = polygons1

	var polygons [][][]float64 = [][][]float64{
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
	for _, pol := range polygons {
		myPoints := make([]*geometry.Point, 0)
		for _, point := range pol {
			myPoint := geometry.NewPoint(point[1], point[0])
			myPoints = append(myPoints, myPoint)
		}
		myPol := geometry.NewPolygon(myPoints)
		arg = append(arg, *myPol)
	}

	sgg := grid.NewSphereGridGraph(2*density, density)
	sgg.DistributeNodes()
	sgg.CreateEdges(arg)

	jsonObj, err := json.Marshal(sgg.GridGraph)
	if err != nil {
		panic(err)
	}

	wErr := os.WriteFile("graph.json", jsonObj, 0644)
	if wErr != nil {
		panic(err)
	}
}
