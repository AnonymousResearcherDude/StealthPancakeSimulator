package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/general"
	"StealthPancakeSimulator/model/parts/utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type LinkInfo struct {
	Count        int
	HopLinkUsage []map[string]int
	LinkUsage    map[string]int
	Paylinks     map[string]int
	NotPaylinks  map[string]int
	File         *os.File
	Writer       *bufio.Writer
}

func InitLinkInfo() *LinkInfo {
	li := LinkInfo{}
	li.LinkUsage = make(map[string]int)
	li.HopLinkUsage = make([]map[string]int, 10)
	for hop := 0; hop < len(li.HopLinkUsage); hop++ {
		li.HopLinkUsage[hop] = make(map[string]int)
	}
	li.Paylinks = make(map[string]int)
	li.NotPaylinks = make(map[string]int)
	li.File = MakeFile("./results/links.txt")
	li.Writer = bufio.NewWriter(li.File)
	LogExpSting(li.Writer)
	return &li
}

func (li *LinkInfo) Reset() {
	li.LinkUsage = make(map[string]int)
	li.HopLinkUsage = make([]map[string]int, 10)
	for hop := 0; hop < len(li.HopLinkUsage); hop++ {
		li.HopLinkUsage[hop] = make(map[string]int)
	}
	li.Paylinks = make(map[string]int)
	li.NotPaylinks = make(map[string]int)
}

func (li *LinkInfo) Close() {
	err := li.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for links output")
	}
	err = li.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/links.txt")
	}
}

func (li *LinkInfo) HopLinkGini() []float64 {
	result := make([]float64, len(li.HopLinkUsage))
	for hop, links := range li.HopLinkUsage {
		list := make([]int, 0, len(links))
		for _, usage := range links {
			list = append(list, usage)
		}
		result[hop] = utils.Gini(list)
	}
	return result
}

func (li *LinkInfo) BucketLinkGini() []float64 {
	bucketlinkusage := make(map[int][]int)
	for link, usage := range li.LinkUsage {
		node1, err := strconv.Atoi(strings.Split(link, "-")[0])
		if err != nil {
			fmt.Println("Error computing BuccketLinkGini: ", err)
		}
		node2, err := strconv.Atoi(strings.Split(link, "-")[1])
		if err != nil {
			fmt.Println("Error computing BuccketLinkGini: ", err)
		}
		lin := config.GetBits() - general.BitLength(node1^node2)
		list, ok := bucketlinkusage[lin]
		if !ok {
			list = make([]int, 0)
		}
		list = append(list, usage)
		bucketlinkusage[lin] = list
	}
	result := make([]float64, len(bucketlinkusage))
	for lin, list := range bucketlinkusage {
		result[lin] = utils.Gini(list)
	}
	return result
}

func (li *LinkInfo) Update(output *Route) {
	li.Count++
	if output.failed() {
		return
	}
	route := output.RouteWithPrices
	payments := output.PaymentsWithPrices
	for h, hop := range route {
		node1 := hop.RequesterNode.ToInt()
		node2 := hop.ProviderNode.ToInt()
		link := fmt.Sprintf("%d-%d", node1, node2)
		li.LinkUsage[link]++
		li.HopLinkUsage[h][link]++
		found := false
		for _, payment := range payments {
			if payment.Payment.FirstNodeId == hop.RequesterNode {
				delete(li.NotPaylinks, link)
				li.Paylinks[link] += hop.Price
				found = true
				break
			}
		}
		if !found {
			li.NotPaylinks[link] += hop.Price
		}
	}
}

func (li *LinkInfo) Log() {
	_, err := li.Writer.WriteString(fmt.Sprintf("\n Current count: %d\n", li.Count))
	if err != nil {
		panic(err)
	}

	_, err = li.Writer.WriteString(fmt.Sprintf("Paylinks number %d and unpayed %d\n", len(li.Paylinks), len(li.NotPaylinks)))
	if err != nil {
		panic(err)
	}

	_, err = li.Writer.WriteString("\n Link usage gini per bucket: ")
	if err != nil {
		panic(err)
	}
	for _, gini := range li.BucketLinkGini() {
		_, err = li.Writer.WriteString(fmt.Sprintf("%.3f, ", gini))
		if err != nil {
			panic(err)
		}
	}

	_, err = li.Writer.WriteString("\n Link usage gini per hop: ")
	if err != nil {
		panic(err)
	}
	for _, gini := range li.HopLinkGini() {
		_, err = li.Writer.WriteString(fmt.Sprintf("%.3f, ", gini))
		if err != nil {
			panic(err)
		}
	}
}
