package output

import (
	"StealthPancakeSimulator/config"
	"bufio"
	"fmt"
	"os"
)

type SuccessInfo struct {
	UniqueCount         int
	Found               int
	FromCache           int
	FromOriginatorCache int
	ThresholdFailed     int
	AccessFailed        int
	File                *os.File
	Writer              *bufio.Writer
}

func InitSuccessInfo() *SuccessInfo {
	si := SuccessInfo{}
	si.File = MakeFile("./results/work.txt")
	si.Writer = bufio.NewWriter(si.File)
	LogExpSting(si.Writer)
	return &si
}

func (si *SuccessInfo) Close() {
	err := si.Writer.Flush()
	if err != nil {
		fmt.Println("Couldn't flush the remaining buffer in the writer for work output")
	}
	err = si.File.Close()
	if err != nil {
		fmt.Println("Couldn't close the file with filepath: ./results/work.txt", err)
	}
}

func (si *SuccessInfo) Reset() {
	si.UniqueCount = 0
	si.Found = 0
	si.FromCache = 0
	si.FromOriginatorCache = 0
	si.AccessFailed = 0
	si.ThresholdFailed = 0
}

func (si *SuccessInfo) Update(output *Route) {
	if output.RetryCount == 0 {
		si.UniqueCount++
	}

	if output.Found {
		si.Found++
	}
	if output.FoundByCaching {
		si.FromCache++
		if output.FoundByOriginatorCache {
			si.FromOriginatorCache++
		}
	}
	if output.AccessFailed {
		si.AccessFailed++
	}
	if output.ThresholdFailed {
		si.ThresholdFailed++
	}
}

func (si *SuccessInfo) Log() {
	total := si.UniqueCount
	foundperc := float64(si.Found) * 100.0 / float64(total)
	_, err := si.Writer.WriteString(fmt.Sprintf("Successfull found: %d, %.2f%%  \n", si.Found, foundperc))
	if err != nil {
		panic(err)
	}

	if config.IsCacheEnabled() {
		cacheperc := float64(si.FromCache) * 100.0 / float64(total)
		_, err = si.Writer.WriteString(fmt.Sprintf("Found from cache: %d, %.2f%%  \n", si.FromCache, cacheperc))
		if err != nil {
			panic(err)
		}
		cacheOriginatorPerc := float64(si.FromOriginatorCache) * 100 / float64(si.FromCache)
		_, err = si.Writer.WriteString(fmt.Sprintf("Found from originator's cache: %d, %.2f%%  \n", si.FromOriginatorCache, cacheOriginatorPerc))
		if err != nil {
			panic(err)
		}
	}

	threshfailperc := float64(si.ThresholdFailed) * 100.0 / float64(total)
	_, err = si.Writer.WriteString(fmt.Sprintf("Threshold failures: %d, %.2f%%  \n", si.ThresholdFailed, threshfailperc))
	if err != nil {
		panic(err)
	}

	accfailperc := float64(si.AccessFailed) * 100.0 / float64(total)
	_, err = si.Writer.WriteString(fmt.Sprintf("Access failures: %d, %.2f%%  \n", si.AccessFailed, accfailperc))
	if err != nil {
		panic(err)
	}
}
