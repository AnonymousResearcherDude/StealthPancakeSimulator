package update

import (
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/types"
)

func Reroute(state *types.State, requestResult types.RequestResult, curEpoch int) int {
	var retryCounter int = 0
	if config.IsRetryWithAnotherPeer() {

		route := requestResult.Route
		chunkId := requestResult.ChunkId
		originatorId := route[0]
		originator := state.Graph.GetNode(originatorId)
		reroute := originator.RerouteStruct.Reroute // reroute = rejected nodes + chunk

		if requestResult.Found {
			if reroute.RejectedNodes != nil {
				if reroute.ChunkId == chunkId { // If chunkId == chunkId --> reset reroute

					originator.RerouteStruct.ResetRerouteAndSaveToHistory(chunkId, curEpoch)
				}
			}

		} else if len(route) > 1 { // Rejection in second hop --> route have at least an originator and a lastHopNode
			lastHopNode := route[len(route)-1]
			if reroute.RejectedNodes == nil {
				reroute = originator.RerouteStruct.AddNewReroute(requestResult.AccessFailed, lastHopNode, chunkId, curEpoch)
			}
			originator.RerouteStruct.AddNodeToRejectedNodes(requestResult.AccessFailed, lastHopNode, chunkId, curEpoch)
		}

		retryCounter = len(reroute.RejectedNodes)

		if len(reroute.RejectedNodes) >= config.GetBinSize() {
			originator.RerouteStruct.ResetRerouteAndSaveToHistory(chunkId, curEpoch)
		}

	}
	return retryCounter
}
