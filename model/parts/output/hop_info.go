package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/utils"
	"bufio"
	"fmt"
	"os"
)

type HopInfo struct {
	HopIncome       map[int]int // income per hop disrespect storer
	HopActionIncome map[int]int // income per hop with storere income on -1
	RouteLength     []int
	File            *os.File
	Writer          *bufio.Writer
}

func InitHopInfo() *HopInfo {
	hinfo := HopInfo{}
	hinfo.HopIncome = make(map[int]int)
	hinfo.RouteLength = make([]int, 0, config.GetIterations())
	hinfo.File = MakeFile("./results/hops.txt")
	hinfo.Writer = bufio.NewWriter(hinfo.File)
	LogExpSting(hinfo.Writer)
	return &hinfo
}

func (hi *HopInfo) Close() {
	err := hi.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for hops output")
	}
	err = hi.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/hops.txt")
	}
}

func (hi *HopInfo) Reset() {
	hi.HopIncome = make(map[int]int)
	hi.RouteLength = make([]int, 0, config.GetIterations())
}

func (hi *HopInfo) CalculateRouteHopIncome() []int {
	result := make([]int, len(hi.HopIncome))

	for i, income := range hi.HopIncome {
		if i == -1 {
			result[len(result)-1] = income
		} else {
			result[i-1] = income
		}
	}

	return result
}

func (hi *HopInfo) CalculateAvgRouteLength() float64 {
	return utils.Mean(hi.RouteLength)
}

func (hi *HopInfo) Update(output *Route) {
	if output.failed() {
		return
	}
	route := output.RouteWithPrices
	if config.GetHopFractionOfRewards() {
		for i, hop := range route {
			if i < len(route)-1 {
				hi.HopIncome[i+1] += hop.Price
			} else {
				// payment to storer
				hi.HopIncome[-1] += hop.Price
			}
			if i > 0 {
				hi.HopIncome[i] -= hop.Price
			}
		}
	}

	if config.GetAverageNumberOfHops() {
		hi.RouteLength = append(hi.RouteLength, len(route))
	}
}

func (hi *HopInfo) Log() {
	if config.GetAverageNumberOfHops() {
		_, err := hi.Writer.WriteString(fmt.Sprintf("Avg route length: %.2f\n", hi.CalculateAvgRouteLength()))
		if err != nil {
			panic(err)
		}
	}

	if config.GetHopFractionOfRewards() {
		routeHopIncome := hi.CalculateRouteHopIncome()
		_, err := hi.Writer.WriteString("RouteHop distribution: \n")
		if err != nil {
			panic(err)
		}
		for hop, income := range routeHopIncome {
			if hop == len(routeHopIncome)-1 {
				_, err = hi.Writer.WriteString(fmt.Sprintf("Hop: storer have income fraction %d\n", income))
			} else {
				_, err = hi.Writer.WriteString(fmt.Sprintf("Hop: %d has income fraction %d\n", hop, income))
			}
			if err != nil {
				panic(err)
			}
		}
	}
}
