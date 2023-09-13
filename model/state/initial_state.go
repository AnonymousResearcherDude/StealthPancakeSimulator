package state

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"StealthPancakeSimulator/model/parts/utils"
	"fmt"
	"math/rand"
)

func MakeInitialState(path string) types.State {
	// Initialize the state
	fmt.Println("start of make initial state")
	rand.Seed(config.GetRandomSeed())
	network := types.Network{}
	network.Load(path)
	graph, err := utils.CreateGraphNetwork(&network)
	if err != nil {
		fmt.Println("create graph network returned an error: ", err)
	}
	//pendingStruct := types.PendingStruct{PendingMap: make(types.PendingMap, 0), PendingMutex: &sync.Mutex{}}
	//rerouteStruct := types.RerouteStruct{RerouteMap: make(types.RerouteMap, 0), RerouteMutex: &sync.Mutex{}}
	//cacheStruct := types.CacheStruct{CacheHits: 0, CacheMap: make(types.CacheMap), CacheMutex: &sync.Mutex{}}

	initialState := types.State{
		Graph:                graph,
		Originators:          utils.CreateDownloadersList(graph),
		RouteLists:           make([]types.RequestResult, 10000),
		UniqueWaitingCounter: 0,
		UniqueRetryCounter:   0,
		OriginatorIndex:      0,
		TimeStep:             0,
		CidsDataLoaded:       utils.GetRealWorkLoadFromFileAndAddToOriginators(graph),
		Epoch:                0,
	}
	return initialState
}
