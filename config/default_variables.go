package config

func SetDefaultConfig() {
	theconfig = getDefaultConfig()
}

func getDefaultConfig() Config {
	return Config{
		BaseOptions: baseOptions{
			Iterations:                      100_000,   // 100_000
			Bits:                            16,        // 16
			NetworkSize:                     10000,     // 10000
			BinSize:                         16,        // 16
			Originators:                     1000,      // 0.01 * NetworkSize
			RefreshRate:                     8,         // 8
			Threshold:                       16,        // 16
			RandomSeed:                      123456789, // 123456789
			MaxProximityOrder:               16,        // 16
			Price:                           1,         // 1
			RequestsPerSecond:               100_000,   // 100_000
			EdgeLock:                        true,      // false
			SameOriginator:                  false,     // false
			IterationMeansUniqueChunk:       false,     // false
			RetryCausesTimeIncrease:         false,     //false
			DebugPrints:                     false,     // false
			DebugInterval:                   1000000,   // 1000000
			NumGoroutines:                   -1,        // -1 means gets overwritten by numCPU
			OutputEnabled:                   false,     // false
			AddressChangeThreshold:          0,         // non-positive means no limit
			OriginatorShuffleProbability:    0.0,       // 0.0
			NonOriginatorShuffleProbability: 0.0,       // 0.0
			ReplicationFactor:               4,
			AdjustableThresholdExponent:     3,
			RealWorkload:                    false,
			OutputOptions: outputOptions{
				MeanRewardPerForward:      false,     // false
				AverageNumberOfHops:       false,     // false
				HopFractionOfTotalRewards: false,     //false
				NegativeIncome:            false,     // false
				IncomeGini:                false,     // false
				IncomeTheil:               false,     // false
				HopIncome:                 false,     // false
				DensenessIncome:           false,     // false
				WorkIncomeSpearman:        false,     // false
				WorkInfo:                  false,     // false
				BucketInfo:                false,     // false
				LinkInfo:                  false,     // false
				ExperimentId:              "default", // default
				Reset:                     false,     // false
				EvaluateInterval:          0,         // 0
			},
		},
		Experiment: experiment{Name: "default"},
		ExperimentOptions: experimentOptions{
			ThresholdEnabled:                  true,  // true
			ReciprocityEnabled:                true,  // true
			ForgivenessEnabled:                true,  // true
			PaymentEnabled:                    false, // false
			MaxPOCheckEnabled:                 false, // false
			OnlyOriginatorPays:                false, // false
			PayOnlyForCurrentRequest:          false, // false
			PayIfOrigPays:                     false, // false
			ForwardersPayForceOriginatorToPay: false, // false
			WaitingEnabled:                    false, // false
			RetryWithAnotherPeer:              false, // false
			PreferredChunks:                   false, // false
			AdjustableThreshold:               false, // false
			CacheIsEnabled:                    false, // false
			CacheSize:                         100000,
			CacheModel: cacheModel{
				Unlimited:    false,
				NonProximity: false,
				LRU:          false,
				LFU:          false,
			},
		},
	}
}
