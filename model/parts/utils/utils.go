package utils

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/general"
	"StealthPancakeSimulator/model/parts/types"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

func SortedKeys(nodeMap map[types.NodeId]*types.Node) []types.NodeId {
	keys := make([]types.NodeId, len(nodeMap))
	i := 0
	for k := range nodeMap {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	if keys[0] < 0 {
		panic("generated network contains a node with an invalid Id")
	}

	return keys
}

func CreateGraphNetwork(net *types.Network) (*types.Graph, error) {
	//fmt.Println("Creating graph network...")
	sortedNodeIds := SortedKeys(net.NodesMap)

	Edges := make(map[types.NodeId]map[types.NodeId]*types.Edge)

	graph := &types.Graph{
		Network: net,
		Edges:   Edges,
	}

	for _, nodeId := range sortedNodeIds {
		graph.Edges[nodeId] = make(map[types.NodeId]*types.Edge)

		node := net.NodesMap[nodeId]
		nodeAdj := node.AdjIds
		for _, adjItems := range nodeAdj {
			for _, otherNodeId := range adjItems {
				threshold := general.BitLength(nodeId.ToInt() ^ otherNodeId.ToInt())
				attrs := types.EdgeAttrs{A2B: 0, LastEpoch: 0, Threshold: threshold}
				err := graph.AddEdge(node.Id, otherNodeId, attrs)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	//fmt.Println("Graph network is created.")
	return graph, nil
}

func GetNewChunkId() types.ChunkId {
	return types.ChunkId(rand.Intn(config.GetAddressRange()-1) + 1)
}

// func GetRealWorkLoadFromFile(currPos int64) *[]types.ChunkId {
// 	if config.GetRealWorkload() {
// 		chunkIds, err := networkdata.GetRealWorkload(currPos)
// 		if err != nil {
// 			fmt.Println("Error reading")
// 		}
// 		return &chunkIds
// 	}

//		return nil
//	}

func GetRealWorkLoadFromFileAndAddToOriginators(g *types.Graph) bool {
	if !config.GetRealWorkload() {
		return false
	}

	allCidData := make([]types.CidStruct, 0)

	file, err := os.Open("network_data/workload.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		cid := row[0]
		frequencyStr := row[1]
		bytesReturnedStr := row[2]

		frequency, err := strconv.Atoi(frequencyStr)
		if err != nil {
			log.Fatal(err)
		}

		bytesReturned, err := strconv.Atoi(bytesReturnedStr)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < frequency; i++ {
			c := types.CidStruct{
				Cid:            cid,
				ReturenedBytes: bytesReturned,
			}
			allCidData = append(allCidData, c)
		}
	}

	rand.Shuffle(len(allCidData), func(i, j int) { allCidData[i], allCidData[j] = allCidData[j], allCidData[i] })

	counter := 0
loop:
	index := 0
	for _, node := range g.NodesMap {
		if node.IsOriginator {
			if counter < len(allCidData) {
				node.ChunksQueueStruct.AddToCidQueue(allCidData[counter])
				counter++
			} else {
				break
			}
		}
		if index == len(g.NodesMap)-1 {
			goto loop
		}
		index++
	}

	return true
}

func GetNewChunkIdFromWorkload(chunkIds *[]types.ChunkId) types.ChunkId {
	if len(*chunkIds) == 0 {
		return -1
	}
	// fmt.Println("chunkid size:", len(*chunkIds))
	poppedChunkId := (*chunkIds)[len(*chunkIds)-1]
	*chunkIds = (*chunkIds)[:len(*chunkIds)-1]
	return poppedChunkId
}

func GetPreferredChunkId() types.ChunkId {
	var chunkId types.ChunkId
	var random float32
	numPreferredChunks := 1
	random = rand.Float32()
	if float32(random) <= 0.8 {
		chunkId = types.ChunkId(rand.Intn(numPreferredChunks))
	} else {
		chunkId = types.ChunkId(rand.Intn(config.GetAddressRange()-numPreferredChunks) + numPreferredChunks)
	}
	return chunkId
}

func FindDistance(first types.NodeId, second types.ChunkId) int {
	return config.GetBits() - general.BitLength(first.ToInt()^second.ToInt())
}

func getProximityChunk(firstNodeId types.NodeId, chunkId types.ChunkId) int {
	retVal := config.GetBits() - general.BitLength(firstNodeId.ToInt()^chunkId.ToInt())
	if retVal <= config.GetMaxProximityOrder() {
		return retVal
	} else {
		return config.GetMaxProximityOrder()
	}
}

func PeerPriceChunk(firstNodeId types.NodeId, chunkId types.ChunkId) int {
	val := (config.GetMaxProximityOrder() - getProximityChunk(firstNodeId, chunkId) + 1) * config.GetPrice()
	return val
}

func CreateDownloadersList(g *types.Graph) []types.NodeId {
	//fmt.Println("Creating downloaders list...")

	downloadersList := make([]types.NodeId, 0)
	counter := 0
	for _, originator := range g.NodesMap {
		downloadersList = append(downloadersList, originator.Id)
		originator.IsOriginator = true
		counter++
		if counter >= config.GetOriginators() {
			break
		}
	}

	//fmt.Println("Downloaders list create...!")
	return downloadersList
}
