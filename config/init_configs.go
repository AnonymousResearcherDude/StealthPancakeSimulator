package config

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

// theconfig This is the current configuration.
var theconfig Config

func InitConfig() {
	config, err := ReadYamlFile("config.yaml")
	if err != nil {
		log.Panicln("Unable to read config file: config.yaml")
	}
	theconfig = config
	ValidateBaseOptions(theconfig.BaseOptions)
	SetExperiment(theconfig)
}

func InitConfigWithPath(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	config, err := ReadYamlFile("config.yaml")
	if err != nil {
		log.Panicln("Unable to read config file: config.yaml")
	}
	theconfig = config
	ValidateBaseOptions(theconfig.BaseOptions)
	SetExperiment(theconfig)
}

func SetExperimentId(id string) {
	theconfig.BaseOptions.OutputOptions.ExperimentId = id
}

func SetMaxPO(maxPO int) {
	theconfig.BaseOptions.MaxProximityOrder = maxPO
}

func ReadYamlFile(filename string) (Config, error) {
	yamlFile, err := os.ReadFile(filename)

	var yamlData Config

	if err != nil {
		log.Printf("yamlFile.Get err :%v ", err)
		return yamlData, err
	}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		log.Panicf("Unmarshal: %v", err)
	}
	return yamlData, nil
}

func SetExperiment(yml Config) {

	switch yml.Experiment.Name {
	case "omega":
		fmt.Println("omega experiment chosen")
		OmegaExperiment()

	case "custom":
		fmt.Println("custom experiment chosen")
		CustomExperiment(yml.ExperimentOptions)

	default:
		fmt.Println("default experiment chosen")
	}
}

func ValidateBaseOptions(configOptions baseOptions) {
	SetNumGoroutines(configOptions.NumGoroutines)
	SetEvaluateInterval(configOptions.OutputOptions.EvaluateInterval)
	SetAddressRange(configOptions.Bits)
	SetStorageDepth(configOptions.ReplicationFactor)
	SetRandomSeed()
}

func SetNumGoroutines(numGoroutines int) {
	if numGoroutines == -1 {
		theconfig.BaseOptions.NumGoroutines = runtime.NumCPU()
	}
}

func SetEvaluateInterval(interval int) {
	if interval < 0 {
		theconfig.BaseOptions.OutputOptions.EvaluateInterval = 0
	}
}

func SetAddressRange(numBits int) {
	if numBits <= 0 {
		theconfig.BaseOptions.AddressRange = int(math.Pow(2, float64(theconfig.BaseOptions.Bits)))
	} else {
		theconfig.BaseOptions.AddressRange = int(math.Pow(2, float64(numBits)))
	}
}

func SetStorageDepth(replicationFactor int) {
	if replicationFactor <= 0 {
		replicationFactor = 4
	}
	depth := 0
	n := GetNetworkSize()
	for n/2 >= replicationFactor {
		n = n / 2
		depth++
	}
	theconfig.BaseOptions.StorageDepth = depth
}

func SetRandomSeed() {
	if theconfig.BaseOptions.RandomSeed == -1 {
		theconfig.BaseOptions.RandomSeed = time.Now().UnixNano()
	}
}

func SetCacheModel(cacheMoelInt int) {
	theconfig.ExperimentOptions.CacheIsEnabled = true
	cacheModel := cacheModel{}
	switch cacheMoelInt {
	case 0:
		cacheModel.Unlimited = true
		cacheModel.NonProximity = false
		cacheModel.LRU = false
		cacheModel.LFU = false
	case 1:
		cacheModel.Unlimited = false
		cacheModel.NonProximity = true
		cacheModel.LRU = false
		cacheModel.LFU = false
	case 2:
		cacheModel.Unlimited = false
		cacheModel.NonProximity = false
		cacheModel.LRU = true
		cacheModel.LFU = false
	case 3:
		cacheModel.Unlimited = false
		cacheModel.NonProximity = false
		cacheModel.LRU = false
		cacheModel.LFU = true
	}
	theconfig.ExperimentOptions.CacheModel = cacheModel
}
