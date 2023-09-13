package output

import (
	"StealthPancakeSimulator/config"
	"sync"
)

func Worker(outputChan chan Route, wg *sync.WaitGroup) {
	defer wg.Done()
	var outputStruct Route
	counter := 0

	loggers := CreateLoggers()
	logInterval := config.GetEvaluateInterval()
	reset := config.DoReset()

	for _, logger := range loggers {
		defer logger.Close()
	}

	for outputStruct = range outputChan {
		counter++

		for _, logger := range loggers {
			logger.Update(&outputStruct)

			if logInterval > 0 && counter%logInterval == 0 {
				logger.Log()
				if reset {
					logger.Reset()
				}
			}
		}
	}
	for _, logger := range loggers {
		logger.Log()
	}
}

func CreateLoggers() []LogResetUpdateCloser {
	loggers := make([]LogResetUpdateCloser, 0)

	successInfo := InitSuccessInfo()
	loggers = append(loggers, successInfo)

	if config.GetAverageNumberOfHops() ||
		config.GetHopFractionOfRewards() ||
		config.GetMeanRewardPerForward() {
		hopInfo := InitHopInfo()
		loggers = append(loggers, hopInfo)
	}

	if config.GetPaymentEnabled() &&
		(config.GetAverageNumberOfHops() ||
			config.GetHopFractionOfRewards() ||
			config.GetMeanRewardPerForward()) {
		hopPaymentInfo := InitHopPaymentInfo()
		loggers = append(loggers, hopPaymentInfo)
	}

	var workIncomeInfo *WorkIncomeInfo
	if config.GetWorkIncomeSpearman() {
		workIncomeInfo := InitWorkIncomeInfo()
		loggers = append(loggers, workIncomeInfo)
	}

	if config.GetNegativeIncome() ||
		config.GetIncomeGini() ||
		config.GetHopIncome() ||
		config.GetIncomeTheil() ||
		config.GetDensnessIncome() {
		if workIncomeInfo == nil {
			loggers = append(loggers, InitIncomeInfo())
		} else {
			loggers = append(loggers, workIncomeInfo.IncomeInfo)
		}
	}

	if config.GetWorkInfo() {
		if workIncomeInfo == nil {
			loggers = append(loggers, InitWorkInfo())
		} else {
			loggers = append(loggers, workIncomeInfo.WorkInfo)
		}
	}

	if config.GetBucketInfo() {
		bucketInfo := InitBucketInfo()
		loggers = append(loggers, bucketInfo)
	}

	if config.GetLinkInfo() {
		linkInfo := InitLinkInfo()
		loggers = append(loggers, linkInfo)
	}

	if config.JustPrintOutPut() {
		outputWriter := InitOutputWriter()
		loggers = append(loggers, outputWriter)
	}
	return loggers
}
