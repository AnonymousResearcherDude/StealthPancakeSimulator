package routing

import (
	"testing"
)

const path = "../utils/testdata/nodes_data_8_10000_0.txt"

// TODO: there should be a complete test for the FindRoute and getNext functions.
func TestFindRoute(t *testing.T) {
	// config.InitConfig()
	// network := &types.Network{}
	// network.Load(path)
	// graph, err := utils.CreateGraphNetwork(network)
	// if err != nil {
	// 	panic(err)
	// }

	// var request types.Request
	// request.OriginatorId = 13825
	// request.ChunkId = 12040
	// route, _, found, _, _, _ := FindRoute(request, graph)
	// expectedRoute := []types.NodeId{13825, 12285, 12046}
	// // expectedPayments := []types.Payment{}
	// expectedFound := true
	// // expectedFlag2 := false
	// // expectedFlag3 := false
	// // expectedFlag4 := false
	// assert.DeepEqual(t, expectedRoute, route)
	// // assert.DeepEqual(t, expectedPayments, paymentList)
	// assert.Equal(t, expectedFound, found)
	// // assert.Equal(t, expectedFlag2, accessFailed)
	// // assert.Equal(t, expectedFlag3, thresholdFailed)
	// // assert.Equal(t, expectedFlag4, foundByCaching)

	// request.OriginatorId = 42372
	// request.ChunkId = 58880
	// route, _, found, _, _, _ = FindRoute(request, graph)
	// expectedRoute = []types.NodeId{42372, 58533, 58944}
	// expectedFound = true
	// assert.DeepEqual(t, expectedRoute, route)
	// assert.Equal(t, expectedFound, found)

	// request.OriginatorId = 26658
	// request.ChunkId = 12258
	// route, _, found, _, _, _ = FindRoute(request, graph)
	// expectedRoute = []types.NodeId{26658, 11042, 12218, 12260}
	// expectedFound = true
	// assert.DeepEqual(t, expectedRoute, route)
	// assert.Equal(t, expectedFound, found)
}

func TestGetNext(t *testing.T) {

}
