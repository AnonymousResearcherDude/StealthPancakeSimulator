package output

import (
	"bufio"
	"fmt"
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/utils"
	"os"
)

type HopPaymentInfo struct {
	HopIncome   map[int]int // income per hop disrespect storer
	RouteLength []int
	FwdIncome   []int
	File        *os.File
	Writer      *bufio.Writer
}

func InitHopPaymentInfo() *HopPaymentInfo {
	hpi := HopPaymentInfo{}
	hpi.HopIncome = make(map[int]int)
	hpi.RouteLength = make([]int, 0, config.GetIterations())
	hpi.File = MakeFile("./results/hopPays.txt")
	hpi.Writer = bufio.NewWriter(hpi.File)
	LogExpSting(hpi.Writer)
	return &hpi
}

func (hpi *HopPaymentInfo) Reset() {
	hpi.HopIncome = make(map[int]int)
	hpi.RouteLength = make([]int, 0, config.GetIterations())
}

func (hpi *HopPaymentInfo) Close() {
	err := hpi.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for hops output")
	}
	err = hpi.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/hopPays.txt")
	}
}

func (hpi *HopPaymentInfo) CalculateRouteHopIncome() []int {
	result := make([]int, len(hpi.HopIncome))

	for i, income := range hpi.HopIncome {
		if i == -1 {
			result[len(result)-1] = income
		} else {
			result[i-1] = income
		}
	}

	return result
}

func (hpi *HopPaymentInfo) CalculateAvgForwardIncome() []int {
	result := make([]int, len(hpi.HopIncome))

	for i, income := range hpi.HopIncome {
		if i == -1 {
			result[len(result)-1] = income
		} else {
			result[i-1] = income
		}
	}

	return result
}

func (hpi *HopPaymentInfo) CalculateMeanStdForwardReward() (mean, std float64) {
	mean = utils.Mean(hpi.FwdIncome)
	std = utils.Stdev(hpi.FwdIncome, mean)
	return mean, std
}

func (hpi *HopPaymentInfo) CalculateAvgRouteLength() float64 {
	return utils.Mean(hpi.RouteLength)
}

func (hpi *HopPaymentInfo) Update(output *Route) {
	if output.failed() {
		return
	}
	payments := output.PaymentsWithPrices
	fwdreward := 0
	for i, hop := range payments {
		if i < len(payments)-1 {
			hpi.HopIncome[i+1] += hop.Price
		} else {
			hpi.HopIncome[i+1] += hop.Price
		}
		if i > 0 {
			hpi.HopIncome[i] -= hop.Price
			fwdreward -= hop.Price
			hpi.FwdIncome = append(hpi.FwdIncome, fwdreward)
		}
		fwdreward = hop.Price
	}

	if config.GetAverageNumberOfHops() {
		hpi.RouteLength = append(hpi.RouteLength, len(payments))
	}
}

func (hpi *HopPaymentInfo) Log() {
	if config.GetAverageNumberOfHops() {
		_, err := hpi.Writer.WriteString(fmt.Sprintf("Avg payment length: %.2f\n", hpi.CalculateAvgRouteLength()))
		if err != nil {
			panic(err)
		}
	}

	if config.GetHopFractionOfRewards() {
		routeHopIncome := hpi.CalculateRouteHopIncome()
		_, err := hpi.Writer.WriteString("RouteHop distribution: \n")
		if err != nil {
			panic(err)
		}
		for hop, income := range routeHopIncome {
			_, err = hpi.Writer.WriteString(fmt.Sprintf("Hop: %d has income fraction %d\n", hop, income))
			if err != nil {
				panic(err)
			}
		}
	}

	if config.GetMeanRewardPerForward() {
		mean, std := hpi.CalculateMeanStdForwardReward()
		_, err := hpi.Writer.WriteString(fmt.Sprintf("Mean and stddevc forward reward: %.5f, %.5f \n", mean, std))
		if err != nil {
			panic(err)
		}
	}
}
