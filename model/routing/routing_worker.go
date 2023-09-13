package routing

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/output"
	"StealthPancakeSimulator/model/parts/types"
	"StealthPancakeSimulator/model/parts/update"
	"fmt"
	"sync"
)

func RoutingWorker(pauseChan chan bool, continueChan chan bool, requestChan chan types.Request, outputChan chan output.Route, globalState *types.State, wg *sync.WaitGroup) {

	defer wg.Done()
	openChannel := true
	var request types.Request
	var requestResult types.RequestResult
	var route []types.NodeId
	var paymentList []types.Payment
	var found bool
	var accessFailed bool
	var thresholdFailed bool
	var foundByCaching bool
	var foundByOriginatorsCache bool

	for {
		select {
		case <-pauseChan:
			continueChan <- true

		case request, openChannel = <-requestChan:
			if !openChannel {
				return
			}

			route, paymentList, found, accessFailed, thresholdFailed, foundByCaching, foundByOriginatorsCache = FindRoute(request, globalState.Graph)

			requestResult = types.RequestResult{
				Route:                  route,
				PaymentList:            paymentList,
				ChunkId:                request.ChunkId,
				Found:                  found,
				AccessFailed:           accessFailed,
				ThresholdFailed:        thresholdFailed,
				FoundByCaching:         foundByCaching,
				FoundByOriginatorCache: foundByOriginatorsCache,
			}

			curTimeStep := request.TimeStep
			output := update.Graph(globalState, requestResult, curTimeStep)

			update.Pending(globalState, requestResult, request.Epoch)
			output.RetryCount = update.Reroute(globalState, requestResult, request.Epoch)
			update.Cache(globalState, requestResult)

			if config.IsOutputEnabled() {
				if config.IsDebugPrints() && config.TimeForDebugPrints(curTimeStep) {
					fmt.Println("outputChan length: ", len(outputChan))
				}
				output.Found = found
				output.ThresholdFailed = thresholdFailed
				output.AccessFailed = accessFailed
				output.FoundByCaching = foundByCaching
				output.FoundByOriginatorCache = foundByOriginatorsCache
				outputChan <- output
			}
		}
	}
}
