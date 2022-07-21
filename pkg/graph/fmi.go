package graph

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func WriteFmi(g Graph, filename string) {

	file, cErr := os.Create(filename)

	if cErr != nil {
		log.Fatal(cErr)
	}
	writer := bufio.NewWriter(file)

	// write number of nodes and number of edges
	writer.WriteString(fmt.Sprintf("%d\n", g.NodeCount()))
	writer.WriteString(fmt.Sprintf("%d\n", g.EdgeCount()))

	// list all nodes structured as "id lat lon"
	for i := 0; i < g.NodeCount(); i++ {
		node := g.GetNode(i)
		writer.WriteString(fmt.Sprintf("%d %f %f\n", i, node.Lat, node.Lon))
	}

	// list all edges structured as "fromId targetId distance"
	for id := 0; id < g.NodeCount(); id++ {
		for _, halfEdge := range g.GetHalfEdgesFrom(id) {
			writer.WriteString(fmt.Sprintf("%d %d %d\n", id, halfEdge.To, halfEdge.Distance))
		}
	}

	writer.Flush()
}

// fmi parse states
const (
	PARSE_NODE_COUNT = iota
	PARSE_EDGE_COUNT = iota
	PARSE_NODES      = iota
	PARSE_EDGES      = iota
)

func NewAdjacencyListFromFmi(filename string) *AdjacencyListGraph {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	numNodes := 0
	numParsedNodes := 0

	alg := AdjacencyListGraph{}
	id2index := make(map[int]int)

	parseState := PARSE_NODE_COUNT
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			// skip empty lines
			continue
		} else if line[0] == '#' {
			// skip comments
			continue
		}

		switch parseState {
		case PARSE_NODE_COUNT:
			if val, err := strconv.Atoi(line); err == nil {
				numNodes = val
				parseState = PARSE_EDGE_COUNT
			}
		case PARSE_EDGE_COUNT:
			parseState = PARSE_NODES
		case PARSE_NODES:
			var id int
			var lat, lon float64
			fmt.Sscanf(line, "%d %f %f", &id, &lat, &lon)
			id2index[id] = alg.NodeCount()
			alg.AddNode(Node{Lon: lon, Lat: lat})
			numParsedNodes++
			if numParsedNodes == numNodes {
				parseState = PARSE_EDGES
			}
		case PARSE_EDGES:
			var from, to, distance int
			fmt.Sscanf(line, "%d %d %d", &from, &to, &distance)
			alg.AddEdge(Edge{From: id2index[from], To: id2index[to], Distance: distance})
		}
	}

	if alg.NodeCount() != numNodes {
		// cannot check edge count because ocean.fmi contains duplicates, which are removed during import
		panic("Invalid parsing result")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &alg
}

func NewAdjacencyArrayFromFmi(filename string) *AdjacencyArrayGraph {
	alg := NewAdjacencyListFromFmi(filename)
	aag := NewAdjacencyArrayFromGraph(alg)
	return aag
}

func WritePartitionedFmi(fg FlaggedGraph, filename string) {

	file, cErr := os.Create(filename)

	if cErr != nil {
		log.Fatal(cErr)
	}
	writer := bufio.NewWriter(file)

	// write number of nodes and number of edges
	writer.WriteString(fmt.Sprintf("%d\n", fg.NodeCount()))
	writer.WriteString(fmt.Sprintf("%d\n", fg.EdgeCount()))

	// list all nodes structured as "id lat lon partition"
	for i := 0; i < fg.NodeCount(); i++ {
		node := fg.GetNode(i)
		writer.WriteString(fmt.Sprintf("%d %f %f %d\n", i, node.Lat, node.Lon, fg.GetPartition(i)))
	}

	// list all edges structured as "fromId targetId weight flag"
	for id := 0; id < fg.NodeCount(); id++ {
		for _, halfEdge := range fg.GetHalfEdgesFrom(id) {
			writer.WriteString(fmt.Sprintf("%d %d %d %d\n", id, halfEdge.To, halfEdge.Weight, halfEdge.Flag))
		}
	}

	writer.Flush()
}

func NewFlaggedAdjacencyListFromFmi(filename string) *FlaggedAdjacencyListGraph {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	numNodes := 0
	numParsedNodes := 0

	falg := FlaggedAdjacencyListGraph{}
	id2index := make(map[int]int)

	parseState := PARSE_NODE_COUNT
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			// skip empty lines
			continue
		} else if line[0] == '#' {
			// skip comments
			continue
		}

		switch parseState {
		case PARSE_NODE_COUNT:
			if val, err := strconv.Atoi(line); err == nil {
				numNodes = val
				parseState = PARSE_EDGE_COUNT
			}
		case PARSE_EDGE_COUNT:
			parseState = PARSE_NODES
		case PARSE_NODES:
			var id int
			var lat, lon float64
			var partitionId uint8
			fmt.Sscanf(line, "%d %f %f %d", &id, &lat, &lon, &partitionId)
			id2index[id] = falg.NodeCount()
			falg.AddNode(Node{Lon: lon, Lat: lat})
			falg.Partitions = append(falg.Partitions, partitionId)
			numParsedNodes++
			if numParsedNodes == numNodes {
				parseState = PARSE_EDGES
			}
		case PARSE_EDGES:
			var from, to, weight int
			var flag uint64
			fmt.Sscanf(line, "%d %d %d %d", &from, &to, &weight, &flag)
			falg.AddHalfEdge(id2index[from], FlaggedHalfEdge{To: id2index[to], Weight: weight, Flag: flag})
		}
	}

	if falg.NodeCount() != numNodes {
		// cannot check edge count because ocean.fmi contains duplicates, which are removed during import
		panic("Invalid parsing result")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &falg
}

func NewFlaggedAdjacencyArrayFromFmi(filename string) *FlaggedAdjacencyArrayGraph {
	falg := NewFlaggedAdjacencyListFromFmi(filename)
	faag := NewFlaggedAdjacencyArrayFromGraph(falg)
	return faag
}
