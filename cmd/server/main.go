/*
 * Ship Routing API
 *
 * Access the global ship routing service via a RESTful API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	"math"
	"net/http"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	io "github.com/dmholtz/graffiti/examples/io"
	g "github.com/dmholtz/graffiti/graph"

	server "github.com/dmholtz/osm-ship-routing/pkg/server/openapi_server"
)

const graphFile = "graphs/ocean_equi_4.fmi"

func main() {

	log.Printf("Loading graph from file %s ...\n", graphFile)

	falg128 := io.NewAdjacencyListFromFmi("graphs/ocean_equi_4_grid_arcflags128.fmi", io.ParsePartGeoPoint, io.ParseLargeFlaggedHalfEdge)
	faag128 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]](falg128)

	log.Printf("Building router ...\n")

	biArcflag128Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128, Transpose: faag128, MaxInitializerValue: math.MaxInt}

	DefaultApiService := &server.DefaultApiService{Graph: faag128, Router: biArcflag128Router}
	DefaultApiController := server.NewDefaultApiController(DefaultApiService)

	router := server.NewRouter(DefaultApiController)
	log.Printf("Server started @ http://localhost:8081")

	log.Fatal(http.ListenAndServe(":8081", router))
}
