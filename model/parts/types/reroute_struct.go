package types

import (
	"StealthPancakeSimulator/model/general"
	"sync"
)

type Reroute struct {
	RejectedNodes []NodeId
	ChunkId       ChunkId
	LastEpoch     int
}

type RerouteStruct struct {
	Reroute      Reroute
	History      map[ChunkId][]NodeId // History of access failure nodes related to a specific chunk
	RerouteMutex *sync.Mutex
}

func (r *RerouteStruct) GetReroute() Reroute {
	return r.Reroute
}

func (r *RerouteStruct) AddNewReroute(accessFail bool, nodeId NodeId, chunkId ChunkId, curEpoch int) Reroute {
	r.RerouteMutex.Lock()
	defer r.RerouteMutex.Unlock()

	var rejectedNodes []NodeId

	if previouslyRejectedNodes := r.History[chunkId]; previouslyRejectedNodes != nil {
		rejectedNodes = previouslyRejectedNodes
	}

	newReroute := Reroute{
		RejectedNodes: rejectedNodes,
		ChunkId:       chunkId,
		//LastEpoch:     curEpoch,
	}

	r.Reroute = newReroute
	return r.Reroute
}

func (r *RerouteStruct) AddNodeToRejectedNodes(accessFail bool, nodeId NodeId, chunkId ChunkId, curEpoch int) {
	r.RerouteMutex.Lock()
	defer r.RerouteMutex.Unlock()

	if historyNodes := r.History[chunkId]; accessFail && !general.Contains(historyNodes, nodeId) {
		if historyNodes != nil {
			r.History[chunkId] = append(historyNodes, nodeId)
		} else {
			r.History[chunkId] = []NodeId{nodeId}
		}
	}
	r.Reroute.RejectedNodes = append(r.Reroute.RejectedNodes, nodeId)

}

func (r *RerouteStruct) ResetRerouteAndSaveToHistory(chunkId ChunkId, curEpoch int) {
	r.RerouteMutex.Lock()
	defer r.RerouteMutex.Unlock()

	// r.History[chunkId] = r.Reroute.RejectedNodes
	r.Reroute = Reroute{}
}

//type RerouteMap map[NodeId]Reroute

//type RerouteStruct struct {
//	RerouteMap           RerouteMap
//	RerouteMutex         *sync.Mutex
//	UniqueRerouteCounter int
//}

//func (r *RerouteStruct) GetRerouteMap() Reroute {
//	r.RerouteMutex.Lock()
//	defer r.RerouteMutex.Unlock()
//	reroute, ok := r.RerouteMap[originator]
//	if ok {
//		return reroute
//	}
//	return Reroute{}
//}

//func (r *RerouteStruct) DeleteReroute(originator NodeId) {
//	r.RerouteMutex.Lock()
//	defer r.RerouteMutex.Unlock()
//	delete(r.RerouteMap, originator)
//}

//func (r *RerouteStruct) AddNewReroute(originator NodeId, nodeId NodeId, chunkId ChunkId, curEpoch int) bool {
//	r.RerouteMutex.Lock()
//	defer r.RerouteMutex.Unlock()
//	_, ok := r.RerouteMap[originator]
//	if !ok {
//		r.RerouteMap[originator] = Reroute{
//			CheckedNodes: []NodeId{nodeId},
//			ChunkId:      chunkId,
//			LastEpoch:    curEpoch,
//		}
//		return true
//	}
//	return false
//}

//func (r *RerouteStruct) AddNodeToReroute(originator NodeId, nodeId NodeId) bool {
//	r.RerouteMutex.Lock()
//	defer r.RerouteMutex.Unlock()
//	routeStruct, ok := r.RerouteMap[originator]
//	if ok {
//		routeStruct.CheckedNodes = append(routeStruct.CheckedNodes, nodeId)
//		r.RerouteMap[originator] = routeStruct
//		return true
//	}
//	return false
//}

//func (r *RerouteStruct) UpdateEpoch(originator NodeId, curEpoch int) int {
//	r.RerouteMutex.Lock()
//	defer r.RerouteMutex.Unlock()
//	routeStruct, ok := r.RerouteMap[originator]
//	if ok {
//		routeStruct.LastEpoch = curEpoch
//		r.RerouteMap[originator] = routeStruct
//		return routeStruct.LastEpoch
//	}
//	return -1
//
//}
