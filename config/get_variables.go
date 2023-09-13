package config

import (
	"fmt"
	"strconv"
)

func GetNumRoutingGoroutines() int {
	num := theconfig.BaseOptions.NumGoroutines
	num-- // for the requestWorker
	if IsOutputEnabled() {
		num-- // for the outputWorker
	}
	if num < 1 {
		if IsOutputEnabled() {
			panic("You need at least 3 goroutines for the requestWorker, routingWorker and outputWorker")
		}
		panic("You need at least 2 goroutines for the requestWorker and routingWorker")
	}
	return num
}

func GetNumGoroutines() int {
	return theconfig.BaseOptions.NumGoroutines
}

// func (c *constant) CreateOriginators(){
// 	c.originators = int(0.001 * float64(c.networkSize))
// }

func IsAdjustableThreshold() bool {
	return theconfig.ExperimentOptions.AdjustableThreshold
}

func GetAdjustableThresholdExponent() int {
	return theconfig.BaseOptions.AdjustableThresholdExponent
}

func GetAddressChangeThreshold() int {
	return theconfig.BaseOptions.AddressChangeThreshold
}

func GetOriginatorShuffleProbability() float32 {
	return theconfig.BaseOptions.OriginatorShuffleProbability
}

func GetNonOriginatorShuffleProbability() float32 {
	return theconfig.BaseOptions.NonOriginatorShuffleProbability
}

func GetRealWorkload() bool {
	return theconfig.BaseOptions.RealWorkload
}

func IsForgivenessEnabled() bool {
	return theconfig.ExperimentOptions.ForgivenessEnabled
}

func IsCacheEnabled() bool {
	return theconfig.ExperimentOptions.CacheIsEnabled
}

func GetCacheSize() int {
	return theconfig.ExperimentOptions.CacheSize
}

func IsPreferredChunksEnabled() bool {
	return theconfig.ExperimentOptions.PreferredChunks
}

func IsRetryWithAnotherPeer() bool {
	return theconfig.ExperimentOptions.RetryWithAnotherPeer
}

func IsForwardersPayForceOriginatorToPay() bool {
	return theconfig.ExperimentOptions.ForwardersPayForceOriginatorToPay
}

func IsPayIfOrigPays() bool {
	return theconfig.ExperimentOptions.PayIfOrigPays
}

func IsPayOnlyForCurrentRequest() bool {
	return theconfig.ExperimentOptions.PayOnlyForCurrentRequest
}

func IsOnlyOriginatorPays() bool {
	return theconfig.ExperimentOptions.OnlyOriginatorPays
}

func IsWaitingEnabled() bool {
	return theconfig.ExperimentOptions.WaitingEnabled
}

func GetMaxPOCheckEnabled() bool {
	return theconfig.ExperimentOptions.MaxPOCheckEnabled
}

func GetThresholdEnabled() bool {
	return theconfig.ExperimentOptions.ThresholdEnabled
}

func GetReciprocityEnabled() bool {
	return theconfig.ExperimentOptions.ReciprocityEnabled
}

func GetPaymentEnabled() bool {
	return theconfig.ExperimentOptions.PaymentEnabled
}

func GetRequestsPerSecond() int {
	return theconfig.BaseOptions.RequestsPerSecond
}

func GetIterations() int {
	return theconfig.BaseOptions.Iterations
}

func GetBits() int {
	return theconfig.BaseOptions.Bits
}

func GetNetworkSize() int {
	return theconfig.BaseOptions.NetworkSize
}

func GetBinSize() int {
	return theconfig.BaseOptions.BinSize
}

func GetAddressRange() int {
	return theconfig.BaseOptions.AddressRange
}

func GetStorageDepth() int {
	return theconfig.BaseOptions.StorageDepth
}

func GetOriginators() int {
	return theconfig.BaseOptions.Originators
}

func GetRefreshRate() int {
	return theconfig.BaseOptions.RefreshRate
}

func GetThreshold() int {
	return theconfig.BaseOptions.Threshold
}

func GetRandomSeed() int64 {
	return theconfig.BaseOptions.RandomSeed
}

func GetMaxProximityOrder() int {
	return theconfig.BaseOptions.MaxProximityOrder
}

func GetPrice() int {
	return theconfig.BaseOptions.Price
}

func GetSameOriginator() bool {
	return theconfig.BaseOptions.SameOriginator
}

func IsEdgeLock() bool {
	return theconfig.BaseOptions.EdgeLock
}

func IsIterationMeansUniqueChunk() bool {
	return theconfig.BaseOptions.IterationMeansUniqueChunk
}

func RetryCausesTimeIncrease() bool {
	return theconfig.BaseOptions.RetryCausesTimeIncrease
}

func IsDebugPrints() bool {
	return theconfig.BaseOptions.DebugPrints
}

func GetDebugInterval() int {
	return theconfig.BaseOptions.DebugInterval
}

func TimeForDebugPrints(timeStep int) bool {
	if IsDebugPrints() {
		return timeStep%GetDebugInterval() == 0
	}
	return false
}

func TimeForNewEpoch(timeStep int) bool {
	return timeStep%GetRequestsPerSecond() == 0
}

func GetReplicationFactor() int {
	return theconfig.BaseOptions.ReplicationFactor
}

func IsOutputEnabled() bool {
	return theconfig.BaseOptions.OutputEnabled
}

func JustPrintOutPut() bool {
	if theconfig.BaseOptions.OutputEnabled &&
		!theconfig.BaseOptions.OutputOptions.MeanRewardPerForward &&
		!theconfig.BaseOptions.OutputOptions.AverageNumberOfHops &&
		!theconfig.BaseOptions.OutputOptions.HopFractionOfTotalRewards &&
		!theconfig.BaseOptions.OutputOptions.NegativeIncome &&
		!theconfig.BaseOptions.OutputOptions.IncomeGini &&
		!theconfig.BaseOptions.OutputOptions.IncomeTheil &&
		!theconfig.BaseOptions.OutputOptions.HopIncome &&
		!theconfig.BaseOptions.OutputOptions.DensenessIncome &&
		!theconfig.BaseOptions.OutputOptions.WorkIncomeSpearman &&
		!theconfig.BaseOptions.OutputOptions.WorkInfo &&
		!theconfig.BaseOptions.OutputOptions.BucketInfo &&
		!theconfig.BaseOptions.OutputOptions.LinkInfo {
		return true
	}
	return false
}

func GetMeanRewardPerForward() bool {
	if theconfig.BaseOptions.OutputEnabled && theconfig.ExperimentOptions.MaxPOCheckEnabled {
		return theconfig.BaseOptions.OutputOptions.MeanRewardPerForward
	}
	return false
}

func GetAverageNumberOfHops() bool {
	if theconfig.BaseOptions.OutputEnabled && theconfig.ExperimentOptions.MaxPOCheckEnabled {
		return theconfig.BaseOptions.OutputOptions.AverageNumberOfHops
	}
	return false
}

func GetHopFractionOfRewards() bool {
	return theconfig.BaseOptions.OutputOptions.HopFractionOfTotalRewards
}

func GetNegativeIncome() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.NegativeIncome
	}
	return false
}

func GetIncomeGini() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.IncomeGini
	}
	return false
}

func GetIncomeTheil() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.IncomeTheil
	}
	return false
}

func GetHopIncome() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.HopIncome
	}
	return false
}

func GetDensnessIncome() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.DensenessIncome
	}
	return false
}

func GetWorkIncomeSpearman() bool {
	if theconfig.ExperimentOptions.PaymentEnabled {
		return theconfig.BaseOptions.OutputOptions.WorkIncomeSpearman
	}
	return false
}

func GetWorkInfo() bool {
	return theconfig.BaseOptions.OutputOptions.WorkInfo
}

func GetBucketInfo() bool {
	return theconfig.BaseOptions.OutputOptions.BucketInfo
}

func GetLinkInfo() bool {
	return theconfig.BaseOptions.OutputOptions.LinkInfo
}

func GetExpeimentId() string {
	return theconfig.BaseOptions.OutputOptions.ExperimentId
}

func DoReset() bool {
	return theconfig.BaseOptions.OutputOptions.Reset
}

func GetEvaluateInterval() (i int) {
	return theconfig.BaseOptions.OutputOptions.EvaluateInterval
}

func GetExperimentString() (exp string) {
	exp = fmt.Sprintf("O%dT%dsS%dk%dTh%dFg%dW%d",
		GetOriginators()*100/GetNetworkSize(),
		GetIterations()/GetRequestsPerSecond(),
		GetIterations(),
		GetBinSize(),
		GetThreshold(),
		GetRefreshRate(),
		GetMaxProximityOrder(),
	)
	if GetPaymentEnabled() {
		exp += "Pay"
	}
	if !GetReciprocityEnabled() {
		exp += "NoRec"
	}
	if IsCacheEnabled() {
		exp += "Cache-"
		exp += strconv.Itoa(GetCacheModel())
		exp = exp + "-" + strconv.Itoa(GetCacheSize())
	}
	if IsPreferredChunksEnabled() {
		exp += "Skew"
	}
	if IsAdjustableThreshold() {
		exp += "FgAdj"
	}
	if GetAddressChangeThreshold() > 0 {
		exp += "AddChangeTh-"
		exp += strconv.Itoa(GetAddressChangeThreshold())
	}
	if GetOriginatorShuffleProbability() > 0 {
		exp += "OrgShProb-"
		s := fmt.Sprintf("%v", GetOriginatorShuffleProbability())
		exp += s
	}
	if GetNonOriginatorShuffleProbability() > 0 {
		exp += "NonOrgShProb-"
		s := fmt.Sprintf("%v", GetNonOriginatorShuffleProbability())
		exp += s
	}

	exp += "-" + GetExpeimentId()
	return exp
}

func GetCacheModel() int {
	if IsCacheEnabled() {
		if theconfig.ExperimentOptions.CacheModel.Unlimited {
			return 0
		}
		if theconfig.ExperimentOptions.CacheModel.NonProximity {
			return 1
		}
		if theconfig.ExperimentOptions.CacheModel.LRU {
			return 2
		}
		if theconfig.ExperimentOptions.CacheModel.LFU {
			return 3
		}
	}

	return -1
}
