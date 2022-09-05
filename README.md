# OSM-Ship-Routing

OSM-Ship-Routing provides the backend of a ship routing service using OSM data. This repository implements the backend server for the [interactive routing frontend](https://github.com/dmholtz/osm-ship-routing-gui).
Moreover, it contains a component for merging OSM coastlines into closed polygons as well as a grid graph builder for the oceans.

Components:

- OSM-Ship-Routing server backend
- Grid Graph Builder
- Coastline Merger

## Setup

Docker is the fastest way to deploy the OSM-Ship-Routing server backend on your system.
The other components require a installation from source.

### Docker

Pull the latest image from [Dockerhub](https://hub.docker.com/repository/docker/dmholtz/osm-ship-routing). Then run the backend server:

`docker run -p 8081:8081 --name osm-server dmholtz:/osm-ship-routing:v2.0.2`

A container with name `osm-server` is started and the routing service is exposed at port 8081.

### Installation from Source

#### Prerequisites

- Go 1.18 or later

#### Merge Coastlines

Extracts coastline segments from a `.pbf` file and merges them to closed coastlines.
The output (i.e. the list of polygons) is either written to the GeoJSON file or to a normal JSON file, which is less verbose than GeoJSON and which we call PolyJSON.

```bash
go mod tidy
go run cmd/merger/main.go
```

#### Graph Builder

Builds a spherical grid graph and implements the point-in-polygon test to check which grid points are in the ocean and will thus become nodes in the graph.
Two types of grids are supported:

- Simple Grid
- Equidistributed Grid

The graph is written to a file in the `.fmi` format.

```bash
go mod tidy
go run cmd/graph-builder/main.go
```

#### OSM-Ship-Routing backend server

Starts a HTTP server at port 8081.
By default, a grid graph with equidistributed nodes on the planet's surface, each having at most four outgoing edges is used.
Routing is done by the `BidirectionalArcFlagRouter` from the [graffiti project](https://github.com/dmholtz/graffiti).

```bash
go mod tidy
go run cmd/server/main.go
```

## Customization

The graph builder supports two grid types and can be customized as follows:

### Simple Grid

Distributes nodes equally along the latidue and longitude axis.

Available Parameters:

- density: The overall number of grid points will be $2 \cdot density^2$.

### Equidistributed Grid

Distributes nodes equally on the planets surface.

Available Parameters:

- nTarget: Number of points to distribute on the surface. The actual number of points may vary slightly.
- meshType: Defines the maximum number of outgoing edges per node. One can choose between four and six neighbors and default value is four neighbors.
