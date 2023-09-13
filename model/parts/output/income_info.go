package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/general"
	"StealthPancakeSimulator/model/parts/utils"
	"bufio"
	"fmt"
	"os"
	"sort"
)

type IncomeInfo struct {
	IncomeMap  map[int]int
	HopMap     map[int][]int
	CostMap    map[int]int
	Requesters map[int]int
	File       *os.File
	Writer     *bufio.Writer
}

func InitIncomeInfo() *IncomeInfo {
	iinfo := IncomeInfo{}
	iinfo.IncomeMap = make(map[int]int)
	iinfo.CostMap = make(map[int]int)
	iinfo.HopMap = make(map[int][]int)
	iinfo.Requesters = make(map[int]int) //This map is currently used to find out who is an originator. This should instead be looked up somewhere else.

	iinfo.File = MakeFile("./results/income.txt")
	iinfo.Writer = bufio.NewWriter(iinfo.File)
	LogExpSting(iinfo.Writer)
	return &iinfo
}

func (ii *IncomeInfo) Reset() {
	ii.IncomeMap = make(map[int]int)
	ii.CostMap = make(map[int]int)
	ii.HopMap = make(map[int][]int)
	ii.Requesters = make(map[int]int) //This map is currently used to find out who is an originator. This should instead be looked up somewhere else.
}

func (ii *IncomeInfo) Close() {
	err := ii.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for output")
	}
	err = ii.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/income.txt")
	}
}

func (o *IncomeInfo) CalculateIncomeFairness() (float64, int, int, float64) {
	size := config.GetNetworkSize()
	vals := make([]int, size)
	i := 0
	println("IncomeMap size", len(o.IncomeMap))
	for _, value := range o.IncomeMap {
		vals[i] = value
		i++
		if i == size {
			break
		}
	}
	min, max := utils.MinAndMax(vals)
	return utils.Gini(vals), min, max, utils.Mean(vals)
}

func (o *IncomeInfo) CalculateNonOIncomeFairness() (float64, int, int, float64) {
	size := config.GetNetworkSize()
	vals := make([]int, 0, size)
	for id, value := range o.IncomeMap {
		if o.Requesters[id] == 0 {
			vals = append(vals, value)
		}
	}
	min, max := utils.MinAndMax(vals)
	return utils.Gini(vals), min, max, utils.Mean(vals)
}

func (o *IncomeInfo) CalculateOriginatorCostFairness() (float64, int, int, float64) {
	size := config.GetNetworkSize()
	vals := make([]int, 0, size)
	for id, value := range o.CostMap {
		if o.Requesters[id] > 0 {
			vals = append(vals, value)
		}
	}
	min, max := utils.MinAndMax(vals)
	return utils.Gini(vals), min, max, utils.Mean(vals) // min, max, average
}

func (o *IncomeInfo) CalculateCostAdjustedOriginatorIncomeFairness() float64 {
	size := config.GetNetworkSize()
	vals := make([]int, 0, size)
	for id, cost := range o.CostMap {
		if o.Requesters[id] > 0 {
			income := o.IncomeMap[id]
			vals = append(vals, income-cost+o.Requesters[id]*16)

		}
	}
	return utils.Gini(vals)
}

func (o *IncomeInfo) CalculateIncomeTheilIndex() float64 {
	size := config.GetNetworkSize()
	vals := make([]int, size)
	i := 0
	for _, value := range o.IncomeMap {
		vals[i] = value
		i++
		if i == size {
			break
		}
	}
	return utils.Theil(vals)
}

func (o *IncomeInfo) CalculateNegativeIncome() (float64, float64) {
	totalNegativeIncomeCounter := 0
	nonOriNegativeIncomeCounter := 0
	for id, value := range o.IncomeMap {
		if value < 0 {
			totalNegativeIncomeCounter += 1
			if o.Requesters[id] == 0 {
				nonOriNegativeIncomeCounter += 1
			}
		}
	}
	totalNegIncome := float64(totalNegativeIncomeCounter) / float64(config.GetNetworkSize())
	nonOriNegIncome := 0.0
	if len(o.Requesters) < config.GetNetworkSize() {
		nonOriNegIncome = float64(nonOriNegativeIncomeCounter) / float64(config.GetNetworkSize()-len(o.Requesters))
	}

	return totalNegIncome, nonOriNegIncome
}

func (ii *IncomeInfo) Update(output *Route) {
	if output.failed() {
		return
	}
	payments := output.PaymentsWithPrices
	for hop, payment := range payments {
		payer := int(payment.Payment.FirstNodeId)
		payee := int(payment.Payment.PayNextId)

		if !(payment.Payment.IsOriginator) {
			ii.IncomeMap[payer] -= payment.Price
		} else {
			ii.Requesters[payer]++
			ii.CostMap[payer] += payment.Price
			if hop != 0 {
				panic("First payment in list is not from originator.")
			}
		}
		ii.IncomeMap[payee] += payment.Price
	}
	route := output.RouteWithPrices
	if !config.GetHopIncome() {
		return
	}
	for hop, path := range route {
		payee := path.ProviderNode.ToInt()

		if _, ok := ii.HopMap[payee]; !ok {
			ii.HopMap[payee] = []int{hop}
		} else {
			ii.HopMap[payee] = append(ii.HopMap[payee], hop)
		}
	}
}

func (ii *IncomeInfo) MaxNonZeroTotal() (max int, nonzero int, total int) {
	for _, value := range ii.IncomeMap {
		total += value
		if value > 0 {
			nonzero++
		}
		if value > max {
			max = value
		}
	}
	return max, nonzero, total
}

// calculate average income for each 10%
func (ii *IncomeInfo) CalculateDistribution() []float64 {
	vals := make([]int, 0, len(ii.IncomeMap))

	total := 0

	for _, value := range ii.IncomeMap {
		vals = append(vals, value)
		total += value
	}
	sort.Slice(vals, func(i2, j int) bool {
		return vals[i2] < vals[j]
	})
	if len(vals) == 0 {
		return nil
	}
	tenpercent := len(vals) / 10
	last := 0
	distribution := make([]float64, 0, 10)
	for last+tenpercent < len(vals) {
		income := 0
		for _, value := range vals[last : last+tenpercent] {
			income += value
		}
		distribution = append(distribution, float64(income)/float64(total))
		last = last + tenpercent
	}
	income := 0
	for _, value := range vals[last:] {
		income += value
	}
	distribution = append(distribution, float64(income)/float64(total))

	return distribution
}

type HopIncome struct {
	Hop    float64
	Income int
	Work   int
}

// calculate average income for each 10% ordered by average hop
func (ii *IncomeInfo) CalculateHopDistribution() (incomeDist, hopDist, workDist []float64) {
	vals := make([]HopIncome, 0, len(ii.IncomeMap))

	for id, hops := range ii.HopMap {
		avghop := utils.Mean(hops)
		vals = append(vals, HopIncome{Hop: avghop, Income: ii.IncomeMap[id], Work: len(hops)})
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Hop < vals[j].Hop
	})
	if len(vals) == 0 {
		return nil, nil, nil
	}
	tenpercent := len(vals) / 10
	last := 0
	incomeDist = make([]float64, 0, 10)
	hopDist = make([]float64, 0, 10)
	workDist = make([]float64, 0, 10)
	for last+tenpercent < len(vals) {
		income := 0
		hops := 0.0
		work := 0
		for _, value := range vals[last : last+tenpercent] {
			income += value.Income
			hops += value.Hop
			work += value.Work
		}
		hopDist = append(hopDist, hops/float64(tenpercent))
		incomeDist = append(incomeDist, float64(income)/float64(tenpercent))
		workDist = append(workDist, float64(work)/float64(tenpercent))
		last = last + tenpercent
	}
	income := 0
	hops := 0.0
	work := 0
	for _, value := range vals[last:] {
		income += value.Income
		hops += value.Hop
		work += value.Work
	}
	incomeDist = append(incomeDist, float64(income)/float64(len(vals[last:])))
	hopDist = append(hopDist, hops/float64(len(vals[last:])))
	workDist = append(workDist, float64(work)/float64(len(vals[last:])))

	return incomeDist, hopDist, workDist
}

func (ii *IncomeInfo) CalculateDensenessDistribution() (mean map[int]float64, std map[int]float64) {
	depth := config.GetStorageDepth()

	regions := make([]int, 0)
	regiondenseness := make(map[int]int)
	regionincome := make(map[int][]int)

	totalincome := 0

	for ida, income := range ii.IncomeMap {
		totalincome += income
		newregion := true
		for regionid := range regions {
			if proximity(ida, regionid) >= depth {
				newregion = false
				regiondenseness[regionid]++
				regionincome[regionid] = append(regionincome[regionid], income)
				break
			}
		}
		if newregion {
			regions = append(regions, ida)
			regiondenseness[ida] = 1
			regionincome[ida] = []int{income}
		}
	}

	densenessincome := make(map[int][]int)
	for r := range regions {
		denseness := regiondenseness[r]
		if _, ok := densenessincome[denseness]; !ok {
			densenessincome[denseness] = make([]int, 0)
		}
		densenessincome[denseness] = append(densenessincome[denseness], regionincome[r]...)
	}
	mean = make(map[int]float64)
	std = make(map[int]float64)
	for d, incomes := range densenessincome {
		mean[d] = utils.Mean(incomes)
		std[d] = utils.Stdev(incomes, mean[d])
	}
	return mean, std
}

func proximity(ida, idb int) int {
	return config.GetBits() - general.BitLength(ida^idb)
}

func (ii *IncomeInfo) AvgHopIncome() (income, count map[int]int) {
	hopincomes := make(map[int][]int, 5)
	for id, hops := range ii.HopMap {
		avghop := int(utils.Mean(hops))
		if _, ok := hopincomes[avghop]; !ok {
			hopincomes[avghop] = []int{ii.IncomeMap[id]}
		} else {
			hopincomes[avghop] = append(hopincomes[avghop], ii.IncomeMap[id])
		}
	}
	avgHopIncome := make(map[int]int, 2)
	avgHopCount := make(map[int]int, 2)
	for hop, income := range hopincomes {
		avgHopIncome[hop] = int(utils.Mean(income))
		avgHopCount[hop] = len(income)
	}
	return avgHopIncome, avgHopCount
}

func (ii *IncomeInfo) Log() {
	if config.GetHopIncome() {
		avgHopIncome, avgHopCount := ii.AvgHopIncome()
		for hop, income := range avgHopIncome {
			_, err := ii.Writer.WriteString(fmt.Sprintf("Hop: %d has income %d and count %d\n", hop, income, avgHopCount[hop]))
			if err != nil {
				panic(err)
			}
		}
		_, err := ii.Writer.WriteString("Distribution is: ")
		if err != nil {
			panic(err)
		}
		for _, avg := range ii.CalculateDistribution() {
			_, err = ii.Writer.WriteString(fmt.Sprintf(", %.2f%%", avg*100))
			if err != nil {
				panic(err)
			}
		}

		_, err = ii.Writer.WriteString("\n")
		if err != nil {
			panic(err)
		}

		_, err = ii.Writer.WriteString("Hop distribution is: ")
		if err != nil {
			panic(err)
		}
		income, hops, work := ii.CalculateHopDistribution()
		for _, avg := range hops {
			_, err = ii.Writer.WriteString(fmt.Sprintf(", %.6f", avg))
			if err != nil {
				panic(err)
			}
		}

		_, err = ii.Writer.WriteString("\n")
		if err != nil {
			panic(err)
		}

		_, err = ii.Writer.WriteString("Hop ordered income distribution is: ")
		if err != nil {
			panic(err)
		}
		for _, avg := range income {
			_, err = ii.Writer.WriteString(fmt.Sprintf(", %.2f", avg))
			if err != nil {
				panic(err)
			}
		}
		_, err = ii.Writer.WriteString("\n")
		if err != nil {
			panic(err)
		}

		_, err = ii.Writer.WriteString("Hop ordered work distribution is: ")
		if err != nil {
			panic(err)
		}
		for _, avg := range work {
			_, err = ii.Writer.WriteString(fmt.Sprintf(", %.2f", avg))
			if err != nil {
				panic(err)
			}
		}
		_, err = ii.Writer.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}

	if config.GetNegativeIncome() {
		negativeIncomeRes, nonOriNegIncome := ii.CalculateNegativeIncome()
		_, err := ii.Writer.WriteString(fmt.Sprintf("Negative income: %f %% \n", negativeIncomeRes*100))
		if err != nil {
			panic(err)
		}

		_, err = ii.Writer.WriteString(fmt.Sprintf("Non originators with negative income: %f %% \n", nonOriNegIncome*100))
		if err != nil {
			panic(err)
		}
	}

	if config.GetIncomeGini() {
		incomeFaireness, min, max, mean := ii.CalculateIncomeFairness()
		_, err := ii.Writer.WriteString(fmt.Sprintf("Income fairness: %f, min income: %d, max income: %d, mean income: %f \n", incomeFaireness, min, max, mean))
		if err != nil {
			panic(err)
		}
		nonOIncomeFairness, min, max, mean := ii.CalculateNonOIncomeFairness()
		_, err = ii.Writer.WriteString(fmt.Sprintf("Non Org Income fairness: %f, min income: %d, max income: %d, mean income: %f \n", nonOIncomeFairness, min, max, mean))
		if err != nil {
			panic(err)
		}

		orgCostFairness, min, max, mean := ii.CalculateOriginatorCostFairness()
		_, err = ii.Writer.WriteString(fmt.Sprintf("Org Cost fairness: %f, min cost: %d, max cost: %d, mean cost: %f \n", orgCostFairness, min, max, mean))
		if err != nil {
			panic(err)
		}
		_, err = ii.Writer.WriteString(fmt.Sprintf("Org Utility fairness: %f \n", ii.CalculateCostAdjustedOriginatorIncomeFairness()))
		if err != nil {
			panic(err)
		}
		max, nonzero, total := ii.MaxNonZeroTotal()
		_, err = ii.Writer.WriteString(fmt.Sprintf("Max, total and nonzero: %d, %d, %d \n", max, total, nonzero))
		if err != nil {
			panic(err)
		}
	}

	if config.GetIncomeTheil() {
		incomeTheilIndex := ii.CalculateIncomeTheilIndex()
		_, err := ii.Writer.WriteString(fmt.Sprintf("Income Theil Index: %f \n", incomeTheilIndex))
		if err != nil {
			panic(err)
		}
	}

	if config.GetDensnessIncome() {
		_, err := ii.Writer.WriteString("Denseness, mean income, std\n")
		if err != nil {
			panic(err)
		}
		means, std := ii.CalculateDensenessDistribution()
		for denseness, mean := range means {
			_, err = ii.Writer.WriteString(fmt.Sprintf("Denseness, %d, %.4f, %.4f\n", denseness, mean, std[denseness]))
			if err != nil {
				panic(err)
			}
		}
	}

}
