package update

import (
	"StealthPancakeSimulator/model/parts/types"
	"sync/atomic"
)

func TimeStep(state *types.State) int {
	curTimeStep := int(atomic.AddInt64(&state.TimeStep, 1))
	return curTimeStep

}
