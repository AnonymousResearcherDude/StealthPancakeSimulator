package types

import (
	"StealthPancakeSimulator/config"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

type Network struct {
	Bits     int
	Bin      int
	NodesMap map[NodeId]*Node
}

type NodeId int

func (n NodeId) ToInt() int {
	return int(n)
}

func (n NodeId) IsNil() bool {
	return n.ToInt() == -1
}

type ChunkId int

func (c ChunkId) ToInt() int {
	return int(c)
}

func (c ChunkId) IsNil() bool {
	return c.ToInt() == 0
}

type jsonFormat struct {
	Bits  int `json:"bits"`
	Bin   int `json:"bin"`
	Nodes []struct {
		Id  int   `json:"id"`
		Adj []int `json:"adj"`
	} `json:"Nodes"`
}

func (network *Network) Load(path string) (int, int, map[NodeId]*Node) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %v", err)
		panic("Unable to open network file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %v", err)
			panic("Unable to close network file")
		}
	}(file)
	decoder := json.NewDecoder(file)

	var test jsonFormat
	err = decoder.Decode(&test)
	if err != nil {
		fmt.Printf("Error decoding file %v: %v", path, err)
		panic("Unable to decode network file")
	}

	network.Bits = test.Bits
	network.Bin = test.Bin
	network.NodesMap = make(map[NodeId]*Node)

	for _, node := range test.Nodes {
		node1 := network.node(NodeId(node.Id))
		for _, adj := range node.Adj {
			node2 := network.node(NodeId(adj))
			node1.add(node2)
		}
	}

	return network.Bits, network.Bin, network.NodesMap
}

func (network *Network) NewNode() *Node {
	// This is a very inefficient way to generate a new node
	// Currently, it adds two-way connections to other nodes

	nodeId := generateIds(1, (1<<network.Bits)-1)[0]
	for _, ok := network.NodesMap[NodeId(nodeId)]; ok; _, ok = network.NodesMap[NodeId(nodeId)] {
		nodeId = generateIds(1, (1<<network.Bits)-1)[0]
	}

	node := network.node(NodeId(nodeId))

	choicenodes := make([]NodeId, 0, len(network.NodesMap))
	for key := range network.NodesMap {
		choicenodes = append(choicenodes, key)
	}
	rand.Shuffle(len(choicenodes), func(i, j int) { choicenodes[i], choicenodes[j] = choicenodes[j], choicenodes[i] })
	for _, adj := range choicenodes {
		added, err := node.add(network.NodesMap[adj])
		if err != nil {
			panic(err)
		}
		if added {
			_, err = network.NodesMap[adj].add(node)
			if err != nil {
				panic(err)
			}
		}
	}

	network.NodesMap[NodeId(nodeId)] = node

	return node
}

func (network *Network) node(nodeId NodeId) *Node {
	if nodeId < 0 || nodeId >= (1<<network.Bits) {
		panic("address out of range")
	}
	res := Node{
		Network: network,
		Id:      nodeId,
		AdjIds:  make([][]NodeId, network.Bits),
		Active:  true,
		OriginatorStruct: OriginatorStruct{
			RequestCount: 0,
		},
		CacheStruct: CacheStruct{
			Size:           config.GetCacheSize(),
			CacheMap:       make(CacheMap),
			CacheFreqMap:   make(CacheFreqMap),
			CacheMutex:     &sync.Mutex{},
			EvictionPolicy: GetCachePolicy(),
		},
		PendingStruct: PendingStruct{
			PendingQueue: nil,
			CurrentIndex: 0,
			PendingMutex: &sync.Mutex{},
		},
		RerouteStruct: RerouteStruct{
			Reroute: Reroute{
				RejectedNodes: nil,
				ChunkId:       0,
				LastEpoch:     0,
			},
			History:      make(map[ChunkId][]NodeId),
			RerouteMutex: &sync.Mutex{},
		},
		AdjLock: sync.RWMutex{},
		ChunksQueueStruct: CidQueueStruct{
			CidQueue:      make([]CidStruct, 0),
			ChunksFromCid: make([]ChunkId, 0),
			CidQueueMutex: &sync.Mutex{},
		},
		IsOriginator: false,
	}
	if len(network.NodesMap) == 0 {
		network.NodesMap = make(map[NodeId]*Node)
	}
	if _, ok := network.NodesMap[nodeId]; !ok {
		network.NodesMap[nodeId] = &res
		return &res
	}
	return network.NodesMap[nodeId]

}

func (network *Network) Generate(count int, random bool) []*Node {
	nodeIds := generateIds(count, (1<<network.Bits)-1)
	if !random {
		nodeIds = generateIdsEven(count, (1<<network.Bits)-1)
	}
	nodes := make([]*Node, 0)
	for _, i := range nodeIds {
		node := network.node(NodeId(i))
		nodes = append(nodes, node)
	}

	for i, node1 := range nodes {
		choicenodes := nodes[i+1:]
		rand.Shuffle(len(choicenodes), func(i, j int) { choicenodes[i], choicenodes[j] = choicenodes[j], choicenodes[i] })
		for _, node2 := range choicenodes {
			added, err := node1.add(node2)
			if err != nil {
				panic(err)
			}
			if added {
				_, err = node2.add(node1)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nodes
}

func (network *Network) Dump(path string) error {
	type NetworkData struct {
		Bits  int `json:"bits"`
		Bin   int `json:"bin"`
		Nodes []struct {
			Id  int   `json:"id"`
			Adj []int `json:"adj"`
		} `json:"nodes"`
	}
	data := NetworkData{network.Bits, network.Bin, make([]struct {
		Id  int   `json:"id"`
		Adj []int `json:"adj"`
	}, 0)}
	for _, node := range network.NodesMap {
		var result []int

		node.AdjLock.RLock()
		for _, list := range node.AdjIds {
			for _, ele := range list {
				result = append(result, int(ele))
			}
		}
		node.AdjLock.RUnlock()

		data.Nodes = append(data.Nodes, struct {
			Id  int   `json:"id"`
			Adj []int `json:"adj"`
		}{Id: int(node.Id), Adj: result})
	}
	file, _ := json.Marshal(data)
	err := os.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func generateIds(totalNumbers int, maxValue int) []int {
	// rand.Seed(time.Now().UnixNano())
	generatedNumbers := make(map[int]bool)
	for len(generatedNumbers) < totalNumbers {
		num := rand.Intn(maxValue-1) + 1
		generatedNumbers[num] = true
	}

	result := make([]int, 0, totalNumbers)
	for num := range generatedNumbers {
		result = append(result, num)
	}
	return result
}

func generateIdsEven(totalNumbers int, maxValue int) []int {
	result := make([]int, 0, totalNumbers)
	step := float64(maxValue) / float64(totalNumbers)
	for id := 0.0; id < float64(maxValue); id += step {
		result = append(result, int(id))
	}
	if len(result) < totalNumbers {
		result = append(result, maxValue)
	}
	return result[:totalNumbers]
}
