package output

import (
	"StealthPancakeSimulator/config"
	"StealthPancakeSimulator/model/parts/types"
	"bufio"
	"fmt"
	"os"
)

type LogResetUpdateCloser interface {
	Log()
	Reset()
	Update(output *Route)
	Close()
}

type Route struct {
	RouteWithPrices        []types.NodePairWithPrice
	PaymentsWithPrices     []types.PaymentWithPrice
	Found                  bool
	AccessFailed           bool
	ThresholdFailed        bool
	FoundByCaching         bool
	RetryCount             int
	FoundByOriginatorCache bool
}

func (o *Route) failed() bool {
	return o.ThresholdFailed || o.AccessFailed
}

func MakeFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func LogExpSting(writer *bufio.Writer) {
	_, err := writer.WriteString(fmt.Sprintf("\n %s \n\n", config.GetExperimentString()))
	if err != nil {
		panic(err)
	}
}
