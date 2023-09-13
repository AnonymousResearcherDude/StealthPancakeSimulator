package types

import "StealthPancakeSimulator/config"

type Request struct {
	TimeStep        int
	Epoch           int
	OriginatorIndex int
	OriginatorId    NodeId
	ChunkId         ChunkId
}

type RequestResult struct {
	Route                  []NodeId
	PaymentList            []Payment
	ChunkId                ChunkId
	Found                  bool
	AccessFailed           bool
	ThresholdFailed        bool
	FoundByCaching         bool
	FoundByOriginatorCache bool
}

type Payment struct {
	FirstNodeId  NodeId
	PayNextId    NodeId
	ChunkId      ChunkId
	IsOriginator bool
}

func (p Payment) IsNil() bool {
	if p.PayNextId == 0 && p.FirstNodeId == 0 && p.ChunkId == 0 {
		return true
	} else {
		return false
	}
}

type Threshold [2]NodeId

//type StateData struct {
//	TimeStep int         `json:"t"`
//	State    StateSubset `json:"s"`
//}

type State struct {
	Graph                *Graph
	Originators          []NodeId
	RouteLists           []RequestResult
	UniqueWaitingCounter int64
	UniqueRetryCounter   int64
	OriginatorIndex      int64
	TimeStep             int64
	CurrPosition         int64 // current position in the dataset file
	ChunkIds             *[]ChunkId
	CidsDataLoaded       bool
	Epoch                int
}

func (s *State) GetOriginatorId(originatorIndex int) NodeId {
	if config.GetAddressChangeThreshold() > 0 {
		nodeId := s.Originators[originatorIndex]
		node := s.Graph.GetNode(nodeId)
		if node == nil {
			panic("Node not found")
		}
		if node.OriginatorStruct.RequestCount > config.GetAddressChangeThreshold() {
			newNode, err := s.Graph.NewNode()
			if err != nil {
				panic(err)
			}
			// The new node is not going to get requests
			newNode.Deactivate()
			newNode.ChunksQueueStruct = node.ChunksQueueStruct
			s.Originators[originatorIndex] = newNode.Id
		}
	}

	return s.Originators[originatorIndex]
}

type NodePairWithPrice struct {
	RequesterNode NodeId
	ProviderNode  NodeId
	Price         int
}

type PaymentWithPrice struct {
	Payment Payment
	Price   int
}
