package main

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	networkdata "StealthPancakeSimulator/network_data"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	binSize := flag.Int("binSize", 16, "Number of nodes in each address table bin, k in Kademlia")

	bits := flag.Int("bits", 16, "address length in bits")
	networkSize := flag.Int("N", 10000, "network size")
	rSeed := flag.Int("rSeed", -1, "random Seed")
	id := flag.String("id", "", "an id")
	count := flag.Int("count", -1, "generate count many networks with ids i0,i1,...")
	random := flag.Bool("random", true, "spread nodes randomly")
	useconfig := flag.Bool("config", false, "use config.yaml to initialize bits, binSize, NetworkSize and randomness")

	flag.Parse()

	if *useconfig {
		config.InitConfig()
		*binSize = config.GetBinSize()
		*bits = config.GetBits()
		*networkSize = config.GetNetworkSize()
		*rSeed = int(config.GetRandomSeed())
	}

	if *rSeed != -1 {
		rand.Seed(int64(*rSeed))
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	println("Parameters:")
	println("binSize:", *binSize, "bits:", *bits, "networkSize:", *networkSize, "rSeed:", *rSeed, "id:", *id, "count:", *count, "random:", *random)

	if *count < 0 {
		filename := "../network_data/" + networkdata.GetNetworkDataName(*bits, *binSize, *networkSize, *id, -1)
		generateAndDump(*bits, *binSize, *networkSize, *random, filename)
	}
	for i := 0; i < *count; i++ {
		filename := "../network_data/" + networkdata.GetNetworkDataName(*bits, *binSize, *networkSize, *id, i)
		generateAndDump(*bits, *binSize, *networkSize, *random, filename)
	}
}

func generateAndDump(bits, binSize, N int, random bool, filename string) {

	network := types.Network{Bits: bits, Bin: binSize}
	network.Generate(N, random)

	err := network.Dump(filename)
	if err != nil {
		panic(fmt.Sprintf("dumping network to file gives error: %v", err))
	}
}
