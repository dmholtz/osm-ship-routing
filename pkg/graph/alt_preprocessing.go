package graph

type LandmarkDistances struct {
	LandmarkId    NodeId
	DistancesFrom []int
	DistancesTo   []int
}

func dijkstraItemsToDistances(dijkstraItems []DjkstraItem) []int {
	distances := make([]int, 0)
	for _, item := range dijkstraItems {
		distances = append(distances, item.Distance)
	}
	return distances
}

func AltPreprocessing(g, gt Graph, landmarks []NodeId) []LandmarkDistances {
	landmarkDistancesCollection := make([]LandmarkDistances, 0)

	for _, landmark := range landmarks {
		// compute distances from landmark l to every node: one-to-all-dijkstra in (forward) graph starting at l
		distancesFrom := dijkstraItemsToDistances(DijkstraOneToAll(g, landmark))

		// compute distances from every node to landmark l: one-to-all-dijsktra in transposed graph starting at l
		distancesTo := dijkstraItemsToDistances(DijkstraOneToAll(gt, landmark))

		landmarkDistances := LandmarkDistances{LandmarkId: landmark, DistancesFrom: distancesFrom, DistancesTo: distancesTo}
		landmarkDistancesCollection = append(landmarkDistancesCollection, landmarkDistances)
	}

	return landmarkDistancesCollection
}
