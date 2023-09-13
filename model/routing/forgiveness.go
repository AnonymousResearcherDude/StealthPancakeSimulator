package routing

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"math"
)

func CheckForgiveness(edgeData types.EdgeAttrs, firstNodeId types.NodeId, secondNodeId types.NodeId, graph *types.Graph, request types.Request) (int, bool) {
	passedTime := request.Epoch - edgeData.LastEpoch

	if passedTime <= 0 {
		return edgeData.A2B, false
	}

	refreshRate := config.GetRefreshRate()
	if config.IsAdjustableThreshold() {
		refreshRate = GetAdjustedRefreshrate(edgeData.Threshold, config.GetThreshold(), config.GetRefreshRate(), config.GetAdjustableThresholdExponent())
	}

	removedDeptAmount := passedTime * refreshRate
	newEdgeData := edgeData
	newEdgeData.A2B -= removedDeptAmount
	if newEdgeData.A2B < 0 {
		newEdgeData.A2B = 0
	}
	newEdgeData.LastEpoch = request.Epoch
	graph.SetEdgeData(firstNodeId, secondNodeId, newEdgeData)

	return newEdgeData.A2B, true
}

func GetAdjustedRefreshrate(adjustedThreshold, threshold, refreshRate, power int) int {
	ratio := float64(adjustedThreshold) / float64(threshold)
	return int(math.Ceil(float64(refreshRate) * math.Pow(ratio, float64(power))))
}
