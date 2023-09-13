package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/utils"
	"bufio"
	"fmt"
	"os"
	"sort"
)

type WorkInfo struct {
	ForwardMap map[int]int
	WorkMap    map[int]int
	Requests   map[int]int
	File       *os.File
	Writer     *bufio.Writer
}

func InitWorkInfo() *WorkInfo {
	winfo := WorkInfo{}
	winfo.ForwardMap = make(map[int]int)
	winfo.WorkMap = make(map[int]int)
	winfo.Requests = make(map[int]int)
	winfo.File = MakeFile("./results/work.txt")
	winfo.Writer = bufio.NewWriter(winfo.File)
	LogExpSting(winfo.Writer)
	return &winfo
}

func (wi *WorkInfo) Close() {
	err := wi.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for work output")
	}
	err = wi.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/work.txt")
	}
}

func (wi *WorkInfo) Reset() {
	wi.ForwardMap = make(map[int]int)
	wi.WorkMap = make(map[int]int)
	wi.Requests = make(map[int]int)
}

func (o *WorkInfo) CalculateWorkFairness() float64 {
	size := config.GetNetworkSize()
	vals := make([]int, size)
	i := 0
	for _, value := range o.WorkMap {
		vals[i] = value
		i++
	}
	return utils.Gini(vals)
}

func (o *WorkInfo) CalculateForwardWorkFairness() float64 {
	size := config.GetNetworkSize()
	vals := make([]int, size)
	i := 0
	for _, value := range o.ForwardMap {
		vals[i] = value
		i++
	}
	return utils.Gini(vals)
}

func (o *WorkInfo) CalculateStorageWorkFairness() float64 {
	size := config.GetNetworkSize()
	vals := make([]int, size)
	i := 0
	for id, value := range o.WorkMap {
		vals[i] = value - o.ForwardMap[id]
		i++
	}
	return utils.Gini(vals)
}

// calculate the maximum work done,
// maximum work done by not originator and
// median work done.
func (o *WorkInfo) CalculateMaxMedianWork() (int, int, int) {
	vals := make([]int, 0, len(o.WorkMap))

	maxfwd := 0

	for id, value := range o.ForwardMap {
		vals = append(vals, value)
		if value > maxfwd && o.Requests[id] == 0 {
			maxfwd = value
		}
	}
	sort.Slice(vals, func(i2, j int) bool {
		return vals[i2] < vals[j]
	})
	if len(vals) == 0 {
		return -1, -1, -1
	}

	return vals[len(vals)-1], maxfwd, vals[len(vals)/2]
}

func (wi *WorkInfo) Update(output *Route) {
	if output.failed() {
		return
	}
	route := output.RouteWithPrices
	for i, hop := range route {
		requester := int(hop.RequesterNode)
		provider := int(hop.ProviderNode)
		if i == 0 {
			wi.Requests[requester]++
		}

		if i != len(route)-1 {
			wi.ForwardMap[provider]++
		}
		wi.WorkMap[provider]++
	}
}

func (wi *WorkInfo) Log() {
	workFairness := wi.CalculateWorkFairness()
	forwardFairness := wi.CalculateForwardWorkFairness()
	max, maxfwd, median := wi.CalculateMaxMedianWork()
	_, err := wi.Writer.WriteString(fmt.Sprintf("Workfairness: %f  \n", workFairness))
	if err != nil {
		panic(err)
	}
	_, err = wi.Writer.WriteString(fmt.Sprintf("Forwardworkfairness: %f  \n", forwardFairness))
	if err != nil {
		panic(err)
	}
	_, err = wi.Writer.WriteString(fmt.Sprintf("Storageworkfairness: %f  \n", wi.CalculateStorageWorkFairness()))
	if err != nil {
		panic(err)
	}
	_, err = wi.Writer.WriteString(fmt.Sprintf("Max, max by non originator, and median work done: %d, %d, %d \n", max, maxfwd, median))
	if err != nil {
		panic(err)
	}
	_, err = wi.Writer.WriteString(fmt.Sprintf("Number of peers doing any work at all: %d \n", len(wi.WorkMap)))
	if err != nil {
		panic(err)
	}
}
