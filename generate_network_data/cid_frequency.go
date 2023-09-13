package main

import (
	"StealthPancakeSimulator/config"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func extract_cid_frequency() {
	if !config.GetRealWorkload() {
		return
	}
	filePath := "data.csv"
	outputFilePath := "output.csv"

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and discard the header line
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	// Initialize a map to store cid frequencies and bytes returned
	cidFreq := make(map[string]int)
	cidBytes := make(map[string]int)

	// Read and process each record
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		cid := record[4]                              // Assuming cid is in the fifth column
		bytesReturned, err := strconv.Atoi(record[2]) // Assuming bytes_returned is in the third column
		if err != nil {
			fmt.Println("Error converting bytes_returned:", err)
			continue
		}

		cidFreq[cid]++
		cidBytes[cid] += bytesReturned
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
	}
	defer outputFile.Close()

	outputWriter := csv.NewWriter(outputFile)
	defer outputWriter.Flush()

	// Write header
	outputWriter.Write([]string{"cid", "frequency", "bytes_returned"})

	// Write data rows
	for cid, freq := range cidFreq {
		if cidBytes[cid] != 0 {
			outputWriter.Write([]string{cid, strconv.Itoa(freq), strconv.Itoa(cidBytes[cid] / cidFreq[cid])})
			// fmt.Printf("CID: %s, Frequency: %d, Avg Bytes Returned: %d\n", cid, freq, cidBytes[cid]/cidFreq[cid])
		}
	}

	fmt.Println("Output written to ", outputFilePath)
}
