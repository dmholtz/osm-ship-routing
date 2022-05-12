package graph

import (
	"bufio"
	"fmt"
	"os"
)

func WriteFmi(g Graph, filename string) {

	file, cErr := os.Create(filename)

	if cErr != nil {
		panic(cErr)
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
	for i := 0; i < g.NodeCount(); i++ {
		for _, edge := range g.GetEdgesFrom(i) {
			writer.WriteString(fmt.Sprintf("%d %d %d\n", edge.From, edge.To, edge.Distance))
		}
	}

	writer.Flush()
}

// TDOO: fmi file -> adjacency array graph
