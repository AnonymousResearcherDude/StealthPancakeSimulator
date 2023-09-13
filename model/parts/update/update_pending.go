package update

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"sync/atomic"
)

func Pending(state *types.State, requestResult types.RequestResult, curEpoch int) int64 {
	var waitingCounter int64 = 0
	if config.IsWaitingEnabled() {
		route := requestResult.Route
		chunkId := requestResult.ChunkId
		originatorId := route[0]
		originator := state.Graph.GetNode(originatorId)
		isNewChunk := false

		if config.IsRetryWithAnotherPeer() {
			if requestResult.ThresholdFailed || requestResult.AccessFailed {
				isNewChunk = originator.PendingStruct.AddPendingChunkId(chunkId, curEpoch)
			} else if requestResult.Found {
				if len(originator.PendingStruct.PendingQueue) > 0 {
					originator.PendingStruct.DeletePendingChunkId(chunkId)
				}
			}

		} else {
			if requestResult.ThresholdFailed {
				isNewChunk = originator.PendingStruct.AddPendingChunkId(chunkId, curEpoch)
			} else if requestResult.Found || requestResult.AccessFailed {
				if len(originator.PendingStruct.PendingQueue) > 0 {
					originator.PendingStruct.DeletePendingChunkId(chunkId)
				}
			}
		}

		if isNewChunk {
			waitingCounter = atomic.AddInt64(&state.UniqueWaitingCounter, 1)
		} else {
			waitingCounter = atomic.LoadInt64(&state.UniqueWaitingCounter)
		}
	}

	return waitingCounter
}
