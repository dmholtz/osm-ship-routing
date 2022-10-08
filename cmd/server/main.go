package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	"github.com/dmholtz/graffiti/examples/io"
	g "github.com/dmholtz/graffiti/graph"

	"github.com/dmholtz/osm-ship-routing/internal/server"

	"github.com/gorilla/mux"
)

const graphFile = "graphs/ocean_equi_4_grid_arcflags128.fmi"

var shipRouterCollection map[string]server.ShipRouter = make(map[string]server.ShipRouter)

func routerId(name string) string {
	id := strings.ReplaceAll(name, " ", "-")
	id = strings.ToLower(id)
	return id
}

// Reports the list of available ship routers
func routers(w http.ResponseWriter, req *http.Request) {
	// allow origin
	w.Header().Add("Access-Control-Allow-Origin", "*")

	type routerDescription struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	routerList := make([]routerDescription, 0)
	for id := range shipRouterCollection {
		name := shipRouterCollection[id].String()
		routerList = append(routerList, routerDescription{Id: id, Name: name})
	}
	json.NewEncoder(w).Encode(routerList)
}

// Computes a route using the respective ship router
func computeRoute(w http.ResponseWriter, req *http.Request) {
	// allow origin
	w.Header().Add("Access-Control-Allow-Origin", "*")

	routerName := mux.Vars(req)["router"]

	// filter out invalid or unavailable routers
	if _, ok := shipRouterCollection[routerName]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	shipRouter := shipRouterCollection[routerName]

	// determine query parameter showSearchSpace
	showSearchSpace := false
	if sss := req.URL.Query().Get("show-search-space"); sss != "" {
		if strings.ToLower(sss) == "true" {
			showSearchSpace = true
		}
	}

	// extract RouteRequest from request body
	var routeRequest server.RouteRequest
	err := json.NewDecoder(req.Body).Decode(&routeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// processing
	log.Printf("Processing RouteRequest %v with searchSpace=%t", routeRequest, showSearchSpace)
	routeResponse := shipRouter.ProcessRequest(routeRequest, showSearchSpace)

	err = json.NewEncoder(w).Encode(routeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	log.Printf("Loading graph from file %s ...\n", graphFile)

	falg128 := io.NewAdjacencyListFromFmi(graphFile, io.ParsePartGeoPoint, io.ParseLargeFlaggedHalfEdge)
	faag128 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]](falg128)

	log.Println("Compute ALT heuristic...")
	landmarks := sp.UniformLandmarks[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]](faag128, 16)
	alt16 := sp.NewAltHeurisitc[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int](faag128, faag128, landmarks[:16])

	log.Printf("Building router ...\n")

	dijkstraRouter := sp.DijkstraRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128}
	biDijkstraRouter := sp.BiDijkstraRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128, Transpose: faag128, MaxInitializerValue: math.MaxInt}
	arcflag128Router := sp.ArcFlagRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128}
	biArcflag128Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128, Transpose: faag128, MaxInitializerValue: math.MaxInt}
	altRouter := sp.AStarRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128, Heuristic: alt16}

	shipRouterCollection[routerId(dijkstraRouter.String())] = server.ShipRouter1[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]{Graph: faag128, Router: dijkstraRouter}
	shipRouterCollection[routerId(biDijkstraRouter.String())] = server.ShipRouter1[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]{Graph: faag128, Router: biDijkstraRouter}
	shipRouterCollection[routerId(arcflag128Router.String())] = server.ShipRouter1[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]{Graph: faag128, Router: arcflag128Router}
	shipRouterCollection[routerId(biArcflag128Router.String())] = server.ShipRouter1[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]{Graph: faag128, Router: biArcflag128Router}
	shipRouterCollection[routerId(altRouter.String())] = server.ShipRouter1[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]]{Graph: faag128, Router: altRouter}

	r := mux.NewRouter()
	r.HandleFunc("/routers", routers).Methods("GET")
	r.HandleFunc("/routers/{router}", computeRoute).Methods("POST")

	server := http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	log.Printf("Server started at http://localhost:8081")
	server.ListenAndServe()
}
