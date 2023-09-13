package workers

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"StealthPancakeSimulator/model/parts/update"
	"StealthPancakeSimulator/model/parts/utils"
	"fmt"
	"sync"
)

func RequestWorker(pauseChan chan bool, continueChan chan bool, requestChan chan types.Request, globalState *types.State, wg *sync.WaitGroup) {

	defer wg.Done()
	requestQueueSize := 10
	counter := 0
	curEpoch := 0
	var chunkId types.ChunkId
	timeStep := 0
	iterations := config.GetIterations()
	numRoutingGoroutines := config.GetNumRoutingGoroutines()

	defer close(requestChan)

	for counter < iterations {
		if len(requestChan) <= requestQueueSize {
			originatorIndex := int(update.OriginatorIndex(globalState, timeStep))
			originatorId := globalState.GetOriginatorId(originatorIndex)
			originator := globalState.Graph.GetNode(originatorId)
			originator.OriginatorStruct.AddRequest()

			// Needed for checks waiting and retry
			chunkId = -1

			if config.IsRetryWithAnotherPeer() {
				rerouteStruct := originator.RerouteStruct

				if len(rerouteStruct.Reroute.RejectedNodes) > 0 {
					chunkId = rerouteStruct.Reroute.ChunkId
				}
			}

			if chunkId == -1 || config.RetryCausesTimeIncrease() {
				// do not count retries towards second load.
				timeStep = update.TimeStep(globalState)

				if config.TimeForNewEpoch(timeStep) {
					curEpoch = update.Epoch(globalState)

					waitForRoutingWorkers(pauseChan, continueChan, numRoutingGoroutines)
					update.Neighbors(globalState)
				}
			}

			if config.IsWaitingEnabled() && chunkId == -1 { // No valid chunkId in reroute
				pendingStruct := originator.PendingStruct

				if pendingStruct.PendingQueue != nil {
					queuedChunk, ok := pendingStruct.GetChunkFromQueue(curEpoch)
					if ok {
						chunkId = queuedChunk.ChunkId
					}
				}
			}

			if config.IsIterationMeansUniqueChunk() {
				if chunkId == -1 { // Only increment the counter chunk is not chosen from waiting or retry
					counter++
				}
			} else {
				counter++ // Increment all iterations
			}

			if config.GetRealWorkload() && chunkId == -1 {
				chunkId = originator.ChunksQueueStruct.GetChunkFromCidQueue()
			}

			if chunkId == -1 && !config.GetRealWorkload() { // No waiting and no retry, and qualify for unique chunk
				chunkId = utils.GetNewChunkId()

				if config.IsPreferredChunksEnabled() {
					chunkId = utils.GetPreferredChunkId()
				}
			}

			if chunkId != -1 { // Should always be true, but just in case
				request := types.Request{
					TimeStep:        timeStep,
					Epoch:           curEpoch,
					OriginatorIndex: originatorIndex,
					OriginatorId:    originatorId,
					ChunkId:         chunkId,
				}
				requestChan <- request
			}

			if config.TimeForDebugPrints(timeStep) {
				fmt.Println("TimeStep is currently:", timeStep)
			}
			if config.TimeForDebugPrints(counter) {
				fmt.Println("Counter is currently:", counter)
			}
		}
	}
}
