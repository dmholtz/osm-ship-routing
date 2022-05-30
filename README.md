# OSM-Ship-Routing

Using is Docker is the fastest way to get the backend of the OSM-Ship-Routing service running.
Beside that, an installation from source gives you access to every component of this project including:

- OSM-Server (backend of the OSM-Ship-Routing service)Coastline Merger
- Dijkstra Benchmarks
- Grid Graph Builder
- Cloastline Merger

## Setup Using Docker

1. Pull the image from [Dockerhub](https://hub.docker.com/repository/docker/dmholtz/osm-ship-routing): `docker pull dmholtz/osm-ship-routing:<TAG>`
2. Start a container: `docker run -p 8081:8081 --name osm-server dmholtz/osm-ship-routing`

Note that `<TAG>` needs to be replaced by a valid tag. Please find all available tags on [Dockerhub](https://hub.docker.com/repository/docker/dmholtz/osm-ship-routing).
Tag `1.0.0` refers to the first submission and tag `latest` referst to the most recent release on Dockerhub.

## Installation from Source

### Prerequisites

The `osm-ship-routing` service is written in [Go](https://go.dev/).
Installing and running requires an installation of the Go Programming Language `v1.18` or later.

### Merge Coastlines

```bash
go run cmd/merger/main.go
```

Extracts coastline segments from a `.pbf` file and merges them to closed coastlines.
The output (i.e. the list of polygons) is either written to the GeoJSON file or to a normal JSON file, which is less verbose than GeoJSON and which we call PolyJSON.

### Graph Builder

```bash
go run cmd/graph-builder/main.go
```

Builds a spherical grid graph and implements the point-in-polygon test to check which grid points are in the ocean and will thus become nodes in the graph.
Two types of grids are supported:

#### Simple Grid

Distributes nodes equally along the latidue and longitude axis.

Available Parameters:

- density: The overall number of grid points will be $2*density^2$.

#### Equidistributed Grid

Distributes nodes equally on the planets surface.

Available Parameters:

- nTarget: Number of points to distribute on the surface. The actual number of points may vary slightly.
- meshType: Defines the maximum number of outgoing edges per node. One can choose between four and six neighbors and default value is four neighbors.

#### Output

The output is written to a file in the `fmi` format.

### Run Dijkstra Benchmarks

```bash
go run cmd/dijkstra/main.go
```

Runs 100 random queries and reports the average elapsed time per query.
By default, a grid graph with equidistributed nodes on the planet's surface, each having at most four outgoing edges is used.

### Start OSM-Server

```bash
go run cmd/server/main.go
```

Starts a HTTP server at port 8081.
By default, a grid graph with equidistributed nodes on the planet's surface, each having at most four outgoing edges is used.
