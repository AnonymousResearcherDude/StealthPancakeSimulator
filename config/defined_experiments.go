package config

// These functions modify the respective fields based changes from default

func OmegaExperiment() {
	theconfig.ExperimentOptions.ThresholdEnabled = false
	theconfig.ExperimentOptions.ForgivenessEnabled = false
	theconfig.ExperimentOptions.MaxPOCheckEnabled = true
}

func CustomExperiment(customExperiment experimentOptions) {
	theconfig.ExperimentOptions = customExperiment
}
