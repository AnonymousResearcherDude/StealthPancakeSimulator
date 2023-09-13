package routing

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/general"
	"StealthPancakeSimulator/model/parts/types"
	"StealthPancakeSimulator/model/parts/utils"
)

// returns the next node in the route, which is the closest node to the route in the previous nodes adjacency list
func getNext(request types.Request, firstNodeId types.NodeId, prevNodePaid bool, graph *types.Graph) (types.NodeId, bool, bool, bool, types.Payment) {
	var nextNodeId types.NodeId = -1
	var payNextId types.NodeId = -1
	var payment types.Payment
	var thresholdFailed bool
	var accessFailed bool
	mainOriginatorId := request.OriginatorId
	chunkId := request.ChunkId
	lastDistance := firstNodeId.ToInt() ^ chunkId.ToInt()

	currDist := lastDistance
	payDist := lastDistance

	bin := config.GetBits() - general.BitLength(lastDistance)

	firstNodeAdjIds := graph.GetNodeAdj(firstNodeId)

	for _, nodeId := range firstNodeAdjIds[bin] {
		dist := nodeId.ToInt() ^ chunkId.ToInt()
		if general.BitLength(dist) >= general.BitLength(lastDistance) {
			panic("Something is wrong. Did try to route to a node that is further from the chunk than myself.")
		}
		if dist >= currDist {
			continue
		}
		if !graph.IsActive(nodeId) {
			continue
		}

		// This means the node is now actively trying to communicate with the other node
		if config.IsEdgeLock() {
			// This is dangerous because it locks all the edges on the route:
			//   Imagine two nodes trying to request each other with distance 2, A -- M -- B
			//   Both of them will lock the first edge on their side (A-M & M-B Respectively), and will request for the other.
			//   Boom! Deadlock.
			graph.LockEdge(firstNodeId, nodeId)
		}

		if !IsThresholdFailed(firstNodeId, nodeId, graph, request) {

			if config.IsRetryWithAnotherPeer() {
				rerouteStruct := graph.GetNode(mainOriginatorId).RerouteStruct
				if rerouteStruct.Reroute.RejectedNodes != nil && general.Contains(rerouteStruct.Reroute.RejectedNodes, nodeId) {
					if config.IsEdgeLock() {
						graph.UnlockEdge(firstNodeId, nodeId)
					}
					continue // skips node that's been part of a failed route before
				}
			}

			thresholdFailed = false

			if config.IsEdgeLock() {
				if !nextNodeId.IsNil() {
					// found new nextNode, release lock on previous found.
					graph.UnlockEdge(firstNodeId, nextNodeId)
				}
				if !payNextId.IsNil() {
					// found new nextNode, without payment, release lock on previous found payNext.
					graph.UnlockEdge(firstNodeId, payNextId)
					payNextId = -1 // IMPORTANT!
				}
			}

			currDist = dist
			nextNodeId = nodeId
		} else {
			thresholdFailed = true

			if config.GetPaymentEnabled() {
				if dist < payDist && nextNodeId.IsNil() {
					if config.IsEdgeLock() && !payNextId.IsNil() {
						graph.UnlockEdge(firstNodeId, payNextId)
					}
					payDist = dist
					payNextId = nodeId
				} else if config.IsEdgeLock() {
					graph.UnlockEdge(firstNodeId, nodeId)
				}
			} else if config.IsEdgeLock() {
				graph.UnlockEdge(firstNodeId, nodeId)
			}
		}
	}

	if !nextNodeId.IsNil() {
		thresholdFailed = false
		accessFailed = false
	} else if !thresholdFailed {
		accessFailed = true
	}

	if config.GetPaymentEnabled() && !payNextId.IsNil() {
		accessFailed = false

		if firstNodeId == mainOriginatorId {
			payment.IsOriginator = true
		}

		if config.IsOnlyOriginatorPays() {
			// Only set payment if the firstNode is the MainOriginator
			if payment.IsOriginator {
				payment.FirstNodeId = firstNodeId
				payment.PayNextId = payNextId
				payment.ChunkId = chunkId
				nextNodeId = payNextId
				thresholdFailed = false
			} else if config.IsEdgeLock() {
				graph.UnlockEdge(firstNodeId, payNextId)
			}
		} else if config.IsPayIfOrigPays() {
			// Pay if the originator pays or if the previous node has paid
			if payment.IsOriginator || prevNodePaid {
				payment.FirstNodeId = firstNodeId
				payment.PayNextId = payNextId
				payment.ChunkId = chunkId
				nextNodeId = payNextId
				thresholdFailed = false
			} else if config.IsEdgeLock() {
				graph.UnlockEdge(firstNodeId, payNextId)
			}
		} else {
			// Always pays
			payment.FirstNodeId = firstNodeId
			payment.PayNextId = payNextId
			payment.ChunkId = chunkId
			nextNodeId = payNextId
			thresholdFailed = false
		}
	}

	prevNodePaid = !payment.IsNil()

	return nextNodeId, thresholdFailed, accessFailed, prevNodePaid, payment
}

func FindRoute(request types.Request, graph *types.Graph) ([]types.NodeId, []types.Payment, bool, bool, bool, bool, bool) {
	chunkId := request.ChunkId
	mainOriginatorId := request.OriginatorId
	curNextNodeId := request.OriginatorId
	route := []types.NodeId{mainOriginatorId}
	found := false
	accessFailed := false
	thresholdFailed := false
	foundByCaching := false
	prevNodePaid := config.IsPayIfOrigPays()
	var payment types.Payment
	var paymentList []types.Payment
	var nextNodeId types.NodeId

	depth := config.GetStorageDepth()

	if config.IsCacheEnabled() {
		mainOriginatorNode := graph.GetNode(mainOriginatorId)
		if mainOriginatorNode.CacheStruct.Contains(chunkId) {
			return nil, nil, true, false, false, true, true
		}
	}

	if utils.FindDistance(mainOriginatorId, chunkId) >= depth {
		found = true
	} else {
		for !(utils.FindDistance(curNextNodeId, chunkId) >= depth) {
			nextNodeId, thresholdFailed, accessFailed, prevNodePaid, payment = getNext(request, curNextNodeId, prevNodePaid, graph)

			if !payment.IsNil() {
				paymentList = append(paymentList, payment)
			}
			if !nextNodeId.IsNil() {
				route = append(route, nextNodeId)
			}
			if !thresholdFailed && !accessFailed {
				if utils.FindDistance(nextNodeId, chunkId) >= depth {
					found = true
					break
				}
				if config.IsCacheEnabled() {
					node := graph.GetNode(nextNodeId)
					if node.CacheStruct.Contains(chunkId) {
						foundByCaching = true
						found = true
						break
					}
				}
				curNextNodeId = nextNodeId
			} else {
				break
			}
		}
	}

	if config.IsForwardersPayForceOriginatorToPay() {
		if !accessFailed && len(paymentList) > 0 {
			newList := make([]types.Payment, 0, len(paymentList))

			for i := 0; i < len(route)-1; i++ {
				newPayment := types.Payment{
					FirstNodeId:  route[i],
					PayNextId:    route[i+1],
					ChunkId:      chunkId,
					IsOriginator: i == 0,
				}
				newList = append(newList, newPayment)

				oldIndex := -1
				for oi, tmp := range paymentList {
					if newPayment.FirstNodeId == tmp.FirstNodeId && newPayment.PayNextId == tmp.PayNextId {
						oldIndex = oi
						break
					}
				}

				if oldIndex > -1 {
					paymentList = append(paymentList[:oldIndex], paymentList[oldIndex+1:]...)
				}
				if len(paymentList) == 0 {
					break
				}
			}

			paymentList = newList
		} else {
			paymentList = []types.Payment{}
		}
	}

	return route, paymentList, found, accessFailed, thresholdFailed, foundByCaching, false // the last one for found fromOriginator's cache
}
