package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/utils"
	"bufio"
	"fmt"
	"os"
)

type WorkIncomeInfo struct {
	WorkInfo   *WorkInfo
	IncomeInfo *IncomeInfo
	File       *os.File
	Writer     *bufio.Writer
}

type SourceRank int

const (
	WorkRank SourceRank = iota
	IncomeRank
)

func InitWorkIncomeInfo() *WorkIncomeInfo {
	wiinfo := WorkIncomeInfo{}
	wiinfo.IncomeInfo = InitIncomeInfo()
	wiinfo.WorkInfo = InitWorkInfo()
	wiinfo.File = MakeFile("./results/work_income.txt")
	wiinfo.Writer = bufio.NewWriter(wiinfo.File)
	LogExpSting(wiinfo.Writer)
	return &wiinfo
}

func (wii *WorkIncomeInfo) Close() {
	err := wii.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for work output")
	}
	err = wii.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/work_income.txt")
	}
}

func (wii *WorkIncomeInfo) Reset() {
	wii.IncomeInfo.Reset()
	wii.WorkInfo.Reset()
}

func (wii *WorkIncomeInfo) CalculateSpearman(sourceRank SourceRank, ratio float64) float64 {
	size := int(float64(config.GetNetworkSize()) * ratio)
	var sourceMap map[int]int
	var otherMap map[int]int
	switch sourceRank {
	case WorkRank:
		sourceMap = wii.WorkInfo.WorkMap
		otherMap = wii.IncomeInfo.IncomeMap
	case IncomeRank:
		sourceMap = wii.IncomeInfo.IncomeMap
		otherMap = wii.WorkInfo.WorkMap
	}
	keys := utils.GetTopKeys(sourceMap, size)
	sourceRanks := utils.GetRanks(sourceMap, keys)
	otherRanks := utils.GetRanks(otherMap, keys)
	return utils.Spearman(sourceRanks, otherRanks)
}

func (wii *WorkIncomeInfo) Update(output *Route) {
	wii.IncomeInfo.Update(output)
	wii.WorkInfo.Update(output)
}

func (wii *WorkIncomeInfo) Log() {
	_, err := wii.Writer.WriteString(fmt.Sprintf("SpearmanAll: %f  \n", wii.CalculateSpearman(WorkRank, 1)))
	if err != nil {
		panic(err)
	}

	_, err = wii.Writer.WriteString(fmt.Sprintf("Spearman10%%Work: %f  \n", wii.CalculateSpearman(WorkRank, 0.1)))
	if err != nil {
		panic(err)
	}
	_, err = wii.Writer.WriteString(fmt.Sprintf("Spearman1%%Work: %f  \n", wii.CalculateSpearman(WorkRank, 0.01)))
	if err != nil {
		panic(err)
	}

	_, err = wii.Writer.WriteString(fmt.Sprintf("Spearman10%%Income: %f  \n", wii.CalculateSpearman(IncomeRank, 0.1)))
	if err != nil {
		panic(err)
	}
	_, err = wii.Writer.WriteString(fmt.Sprintf("Spearman1%%Income: %f  \n", wii.CalculateSpearman(IncomeRank, 0.01)))
	if err != nil {
		panic(err)
	}
}
