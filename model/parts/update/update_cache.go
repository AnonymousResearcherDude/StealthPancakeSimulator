package update

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
)

func Cache(state *types.State, requestResult types.RequestResult) bool {
	if config.IsCacheEnabled() {
		route := requestResult.Route
		chunkId := requestResult.ChunkId

		if requestResult.Found {
			for i, nodeId := range route {
				if i == len(route)-1 && !requestResult.FoundByCaching {
					// do not cache chunks you are responsible for
					continue
				}
				// if utils.PeerPriceChunk(nodeId, chunkId) < config.GetMaxProximityOrder()/2 {
				// 	continue
				// }
				node := state.Graph.GetNode(nodeId)
				node.CacheStruct.AddToCache(chunkId, nodeId)
			}
			return true
		}

	}
	return false
}
