package types

import (
	"StealthPancakeSimulator/model/general"
	"fmt"
	"sync"
)

// Graph structure, node Ids in array and edges in map
type Graph struct {
	*Network
	CurState State
	Edges    map[NodeId]map[NodeId]*Edge
	rwMutex  sync.RWMutex
}

// Edge that connects to NodesMap with attributes about the connection
type Edge struct {
	FromNodeId NodeId
	ToNodeId   NodeId
	Attrs      EdgeAttrs
	Mutex      *sync.Mutex
}

// EdgeAttrs Edge attributes structure,
// "a2b" show how much this node asked from other node,
// "lastEpoch" is the epoch where it was last forgiven.
// "threshold" is for the adjustable threshold limit.
type EdgeAttrs struct {
	A2B       int
	LastEpoch int
	Threshold int
}

func (g *Graph) GetNodeAdj(nodeId NodeId) [][]NodeId {
	n := g.GetNode(nodeId)
	if n == nil {
		panic(fmt.Sprintf("Node %d does not exist!", nodeId))
	}
	n.AdjLock.RLock()
	defer n.AdjLock.RUnlock()
	return n.AdjIds
}

// AddEdge will add an edge from a node to a node
func (g *Graph) AddEdge(fromNodeId NodeId, toNodeId NodeId, attrs EdgeAttrs) error {
	toNode := g.NodesMap[toNodeId]
	fromNode := g.NodesMap[fromNodeId]
	if toNode == nil || fromNode == nil {
		return fmt.Errorf("not a valid edge from %d ---> %d", fromNode.Id, toNode.Id)
	} else if g.unsafeEdgeExists(fromNodeId, toNodeId) {
		return fmt.Errorf("edge from node %d ---> %d already exists", fromNodeId, toNodeId)
	} else {
		mutex := &sync.Mutex{}
		if g.unsafeEdgeExists(toNodeId, fromNodeId) {
			mutex = g.Edges[toNodeId][fromNodeId].Mutex
		}
		newEdge := &Edge{FromNodeId: fromNodeId, ToNodeId: toNodeId, Attrs: attrs, Mutex: mutex}
		g.Edges[fromNodeId][toNodeId] = newEdge
		return nil
	}
}

func (g *Graph) NewNode() (*Node, error) {
	g.rwMutex.Lock()
	node := g.Network.NewNode()
	node.Deactivate()
	g.Edges[node.Id] = make(map[NodeId]*Edge)
	defer g.rwMutex.Unlock()

	nodeAdj := node.AdjIds
	for _, adjItems := range nodeAdj {
		for _, otherNodeId := range adjItems {
			threshold := general.BitLength(node.Id.ToInt() ^ otherNodeId.ToInt())
			attrs := EdgeAttrs{A2B: 0, LastEpoch: 0, Threshold: threshold}
			err := g.AddEdge(node.Id, otherNodeId, attrs)
			if err != nil {
				return nil, err
			}
			err = g.AddEdge(otherNodeId, node.Id, attrs)
			if err != nil {
				return nil, err
			}
		}
	}
	node.Activate()

	return node, nil
}

func (g *Graph) LockEdge(nodeA NodeId, nodeB NodeId) {
	edge := g.GetEdge(nodeA, nodeB)
	if edge == nil {
		panic(fmt.Sprintf("Trying to lock edge %d-%d that does not exist!", nodeA, nodeB))
	}
	edge.Mutex.Lock()
}

func (g *Graph) UnlockEdge(nodeA NodeId, nodeB NodeId) {
	// fmt.Printf("\n UnLockEdge: %d-%d", nodeA, nodeB)
	if !g.EdgeExists(nodeA, nodeB) {
		panic(fmt.Sprintf("Trying to unlock edge %d-%d that does not exist!", nodeA, nodeB))
	}
	edge := g.GetEdge(nodeA, nodeB)
	edge.Mutex.Unlock()
}

func (g *Graph) GetEdge(fromNodeId NodeId, toNodeId NodeId) *Edge {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()
	if _, ok := g.Edges[fromNodeId][toNodeId]; ok {
		return g.Edges[fromNodeId][toNodeId]
	}

	err := g.AddEdge(fromNodeId, toNodeId, EdgeAttrs{})
	if err != nil {
		panic(err.Error())
	}

	return g.Edges[fromNodeId][toNodeId]
}

func (g *Graph) GetEdgeData(fromNodeId NodeId, toNodeId NodeId) EdgeAttrs {
	if g.EdgeExists(fromNodeId, toNodeId) {
		return g.GetEdge(fromNodeId, toNodeId).Attrs
	}
	return EdgeAttrs{}
}

func (g *Graph) EdgeExists(fromNodeId NodeId, toNodeId NodeId) bool {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()
	if _, ok := g.Edges[fromNodeId][toNodeId]; ok {
		return true
	}
	return false
}

func (g *Graph) unsafeEdgeExists(fromNodeId NodeId, toNodeId NodeId) bool {
	if _, ok := g.Edges[fromNodeId][toNodeId]; ok {
		return true
	}
	return false
}

func (g *Graph) SetEdgeData(fromNodeId NodeId, toNodeId NodeId, edgeAttrs EdgeAttrs) bool {
	if g.EdgeExists(fromNodeId, toNodeId) {
		g.Edges[fromNodeId][toNodeId].Attrs = edgeAttrs
		return true
	}
	return false
}

// GetNode getNode will return a node point if exists or return nil
func (g *Graph) GetNode(nodeId NodeId) *Node {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()
	node, ok := g.NodesMap[nodeId]
	if ok {
		return node
	}
	return nil
}

func (g *Graph) IsActive(nodeId NodeId) bool {
	node := g.GetNode(nodeId)
	if node == nil {
		return false
	}
	return node.Active
}

func (g *Graph) Print() {
	for _, v := range g.NodesMap {
		fmt.Printf("%d : ", v.Id)
		for _, i := range v.AdjIds {
			for _, v := range i {
				fmt.Printf("%d ", v)
			}
		}
		fmt.Println()
	}
}
