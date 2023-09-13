package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Request struct {
	Freq   int
	Chunks int
}

func generateUniqueRandomNumbers(count int, existingNumbers map[int]bool) []int {
	uniqueNumbers := make([]int, 0, count)
	rand.Seed(time.Now().UnixNano())

	for len(uniqueNumbers) < count {
		num := rand.Intn(1 << 26)
		if !existingNumbers[num] {
			existingNumbers[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}

	return uniqueNumbers
}

type Record struct {
	Number int
	Freq   int
}

func generateUniqueNumbers() {
	var allReq []Request
	var allRecords []Record

	file, err := os.Open("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	binaryFile, err := os.Create("output.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer binaryFile.Close()

	existingNumbers := make(map[int]bool)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		bytesReturned, err := strconv.Atoi(record[2])
		if err != nil {
			fmt.Println("Error converting bytes_returned:", err)
			continue
		}

		freq, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Error converting frequency:", err)
			continue
		}

		chunks := bytesReturned / 4096
		req := Request{Freq: freq, Chunks: chunks}
		allReq = append(allReq, req)
	}

	for _, req := range allReq {
		randomNumbers := generateUniqueRandomNumbers(req.Chunks, existingNumbers)
		for _, num := range randomNumbers {
			record := Record{Number: num, Freq: req.Freq}
			allRecords = append(allRecords, record)
		}
	}

	for _, record := range allRecords {
		err := binary.Write(binaryFile, binary.LittleEndian, int32(record.Freq))
		if err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
		err = binary.Write(binaryFile, binary.LittleEndian, int32(record.Number))
		if err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}
}
