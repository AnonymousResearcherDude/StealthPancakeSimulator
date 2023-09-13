package types

import (
	"errors"
	"math"

	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/general"
	"math/rand"
	"sync"
)

type Node struct {
	Network           *Network
	Id                NodeId
	Active            bool
	AdjIds            [][]NodeId
	OriginatorStruct  OriginatorStruct
	CacheStruct       CacheStruct
	PendingStruct     PendingStruct
	RerouteStruct     RerouteStruct
	AdjLock           sync.RWMutex
	ChunksQueueStruct CidQueueStruct
	IsOriginator      bool
}

// Adds a one-way connection from node to other
func (node *Node) add(other *Node) (bool, error) {
	if node.Network == nil || node.Network != other.Network {
		return false, errors.New("trying to add nodes with different networks")
	}
	if node == other {
		return false, nil
	}
	if !other.Active {
		return false, nil
	}

	node.AdjLock.Lock()
	defer node.AdjLock.Unlock()

	if node.AdjIds == nil {
		node.AdjIds = make([][]NodeId, node.Network.Bits)
	}
	bit := node.Network.Bits - general.BitLength(node.Id.ToInt()^other.Id.ToInt())
	if bit < 0 || bit >= node.Network.Bits {
		return false, errors.New("nodes have distance outside XOR metric")
	}
	if len(node.AdjIds[bit]) < node.Network.Bin && !general.Contains(node.AdjIds[bit], other.Id) {
		node.AdjIds[bit] = append(node.AdjIds[bit], other.Id)
		return true, nil
	}
	return false, nil
}

func (node *Node) UpdateNeighbors() {
	node.AdjLock.Lock()
	defer node.AdjLock.Unlock()

	candidateNeighbors := make([][]NodeId, node.Network.Bits)
	numConsidredNeighbors := int(math.Log2(float64(node.Network.Bin + 4))) // 4 is an arbitrary smoothing factor
	for l, adjIds := range node.AdjIds {
		shuffledAdjIds := getRandomElements(adjIds, numConsidredNeighbors)
		for _, adjId := range shuffledAdjIds {
			if !general.Contains(candidateNeighbors[l], adjId) {
				candidateNeighbors[l] = append(candidateNeighbors[l], adjId)
			}
			adj := node.Network.NodesMap[adjId]
			adj.AdjLock.RLock()
			for _, adjAdjIds := range adj.AdjIds {
				shuffledAdjAdjIds := getRandomElements(adjAdjIds, numConsidredNeighbors)
				for _, adjAdjId := range shuffledAdjAdjIds {
					bin := config.GetBits() - general.BitLength(node.Id.ToInt()^adjAdjId.ToInt())
					if adjAdjId != node.Id && !general.Contains(candidateNeighbors[bin], adjAdjId) {
						candidateNeighbors[bin] = append(candidateNeighbors[bin], adjAdjId)
					}
				}
			}
			adj.AdjLock.RUnlock()
		}
	}

	for d := 0; d < node.Network.Bits; d++ {
		rand.Shuffle(len(candidateNeighbors[d]), func(i, j int) {
			candidateNeighbors[d][i], candidateNeighbors[d][j] = candidateNeighbors[d][j], candidateNeighbors[d][i]
		})
		if len(candidateNeighbors[d]) > node.Network.Bin {
			node.AdjIds[d] = candidateNeighbors[d][:node.Network.Bin]
		} else {
			node.AdjIds[d] = candidateNeighbors[d]
		}
	}
}

func (node *Node) IsNil() bool {
	return node.Id == 0
}

func (node *Node) Deactivate() {
	node.Active = false
}

func (node *Node) Activate() {
	node.Active = true
}

func getRandomElements(slice []NodeId, num int) []NodeId {
	if len(slice) <= num {
		return slice
	}

	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice[:num]
}
