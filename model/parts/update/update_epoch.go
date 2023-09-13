package update

import "StealthPancakeSimulator/model/parts/types"

func Epoch(state *types.State) int {
	state.Epoch++
	return state.Epoch
}
