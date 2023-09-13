package networkdata

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"

	"gonum.org/v1/gonum/stat"
)

func statisticalZipf() {
	file, err := os.Open("workload.bin")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	freqMap := make(map[int]int)

	for {
		var frequency uint32
		var id uint32
		err := binary.Read(file, binary.LittleEndian, &frequency)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading frequency: ", err)
			return
		}

		err = binary.Read(file, binary.LittleEndian, &id)
		if err != nil {
			fmt.Println("Error reading number: ", err)
			return
		}
		freqMap[int(id)] = int(frequency)
	}

	frequencies := make([]float64, 0, len(freqMap))
	for _, freq := range freqMap {
		frequencies = append(frequencies, float64(freq))
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(frequencies)))

	totalFreq := 0.0
	for _, freq := range frequencies {
		totalFreq += freq
	}

	expectedFreq := make([]float64, len(frequencies))
	for i := range expectedFreq {
		expectedFreq[i] = float64(totalFreq) / float64(i+1)
	}

	pValue := stat.ChiSquare(frequencies, expectedFreq)
	degreesOfFreedom := len(frequencies) - 1

	fmt.Printf("Degrees of freedom: %d\n", degreesOfFreedom)
	fmt.Printf("p-value: %f\n", pValue)

	significanceLevel := 0.05
	if pValue < significanceLevel {
		fmt.Println("Reject the null hypothesis: Data does not follow Zipf's distribution.")
	} else {
		fmt.Println("Fail to reject the null hypothesis: Data follows Zipf's distribution.")
	}
}
