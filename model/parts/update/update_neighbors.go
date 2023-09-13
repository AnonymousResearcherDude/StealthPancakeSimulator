package update

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"math/rand"
)

func Neighbors(globalState *types.State) bool {
	// Update neighbors with probability p
	if config.GetOriginatorShuffleProbability() <= 0 && config.GetNonOriginatorShuffleProbability() <= 0 {
		return true
	}
	for _, node := range globalState.Graph.NodesMap {
		if node.OriginatorStruct.RequestCount > 0 {
			// Originators
			if rand.Float32() < config.GetOriginatorShuffleProbability() {
				node.UpdateNeighbors()
			}
		} else {
			// Non-originators
			if rand.Float32() < config.GetNonOriginatorShuffleProbability() {
				node.UpdateNeighbors()
			}
		}
	}

	return true
}
