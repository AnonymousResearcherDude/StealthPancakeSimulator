package utils

import (
	"fmt"
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"math"
	"math/rand"
	"testing"

	"gotest.tools/assert"
)

const path = "testdata/nodes_data_8_10000_0.txt"

func TestCreateGraphNetwork(t *testing.T) {
	// fileName := "input_test.txt"
	network := &types.Network{}
	network.Load(path)
	graph, err := CreateGraphNetwork(network)

	edge := graph.GetEdge(49584, 0)
	graph.LockEdge(49584, 0)
	graph.UnlockEdge(49584, 0)
	node := graph.GetNode(0)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(graph.NodesMap), 10000)
	assert.Equal(t, len(graph.Edges), 10000)
	assert.Check(t, *edge != types.Edge{})
	assert.Check(t, node != nil)
	edge = graph.GetEdge(27481, 46283)
	assert.Equal(t, edge.Attrs.A2B, 0)
	for node := range graph.Edges {
		for _, edge := range graph.Edges[node] {
			assert.Equal(t, edge.Attrs.A2B, 0)
		}
	}
}

func TestBinSize(t *testing.T) {
	// fileName := "input_test.txt"
	network := &types.Network{}
	network.Load(path)
	graph, err := CreateGraphNetwork(network)

	nodeId := types.NodeId(47391)
	bins := graph.GetNodeAdj(nodeId)
	assert.Equal(t, err, nil)
	assert.Equal(t, network.Bin, 8)
	assert.Equal(t, graph.Bin, 8)
	assert.Equal(t, len(bins[0]), 8)
	assert.Equal(t, len(bins[1]), 7)
}

// TODO: not working right now
func TestCreateDowloaderList(t *testing.T) {
	// Create a network
	network := &types.Network{}
	// Load data to network
	network.Load(path)
	// Creates graph
	graph, _ := CreateGraphNetwork(network)
	// Get number of originators used in the func
	config.SetDefaultConfig()

	c := config.GetOriginators()

	// Create a list of downloaders
	l := CreateDownloadersList(graph)

	// Check if the length of the list is equal to the number of originators specified
	assert.Equal(t, len(l), c)
}

func TestDistributionRespNodeswithStorageDepth(t *testing.T) {
	network := &types.Network{}
	network.Load(path)
	addrRange := math.Pow(2, float64(network.Bits))
	sortedNodeIds := SortedKeys(network.NodesMap)
	fmt.Printf("First and last node %d, %d, length %d \n", sortedNodeIds[0], sortedNodeIds[len(sortedNodeIds)-1], len(sortedNodeIds))
	n := len(sortedNodeIds)
	depth := 0
	for n > 8 {
		n = n / 2
		depth++
	}
	hits := make([]int, 100)
	for i := 0; i < 100; i++ {
		chunkId := types.ChunkId(rand.Intn(int(addrRange)-1) + 1)
		for _, id := range sortedNodeIds {
			if getProximityChunk(id, chunkId) >= depth {
				hits[i]++
			}
		}
	}

	fmt.Printf("Depth %d gives responsible groups %v", depth, hits)
}

func TestGini(t *testing.T) {
	values := []int{4, 0, 0, 0}

	assert.Equal(t, Mean(values), float64(1))
	assert.Equal(t, Gini(values), 0.75)
}

func TestTheil(t *testing.T) {
	values := []int{7, 2, 1, 2}

	println(Theil(values))

	assert.Assert(t, 0.267 < Theil(values))
	assert.Assert(t, Theil(values) < 0.268)
}

func TestSpearman(t *testing.T) {
	rankSet1 := []int{1, 5, 2, 4, 3}
	rankSet2 := []int{3, 2, 4, 1, 5}

	assert.Equal(t, Spearman(rankSet1, rankSet2), -0.5)
}

func TestGetTopKeys(t *testing.T) {
	m := map[int]int{
		1: 8,
		2: 2,
		3: 4,
		4: 6,
	}

	keys := GetTopKeys(m, 3)

	assert.Equal(t, len(keys), 3)
	assert.Equal(t, keys[0], 1)
	assert.Equal(t, keys[1], 4)
	assert.Equal(t, keys[2], 3)
}

func TestGetRanks(t *testing.T) {
	m := map[int]int{
		10: 8,
		20: 2,
		30: 4,
		40: 6,
	}

	ranks := GetRanks(m, []int{20, 10, 30})

	assert.Equal(t, len(ranks), 3)
	assert.Equal(t, ranks[0], 2)
	assert.Equal(t, ranks[1], 0)
	assert.Equal(t, ranks[2], 1)
}
