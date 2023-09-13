package utils

import (
	"sort"
)

func Spearman(rank1 []int, rank2 []int) float64 {
	if len(rank1) != len(rank2) {
		panic("ranks must be of equal length")
	}
	n := len(rank1)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += float64(rank1[i]-rank2[i]) * float64(rank1[i]-rank2[i])
	}
	return 1 - 6*sum/float64(n*(n*n-1))
}

func GetTopKeys(mp map[int]int, size int) []int {
	var kvPairs []struct {
		key   int
		value int
	}

	for k, v := range mp {
		kvPairs = append(kvPairs, struct {
			key   int
			value int
		}{k, v})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].value > kvPairs[j].value
	})

	var result []int
	for i := 0; i < size && i < len(kvPairs); i++ {
		result = append(result, kvPairs[i].key)
	}

	return result
}

func GetRanks(mp map[int]int, keys []int) []int {
	var kvPairs []struct {
		key   int
		value int
	}

	for _, k := range keys {
		// If a key does not exist in the map, it is assumed to have a value of 0
		kvPairs = append(kvPairs, struct {
			key   int
			value int
		}{k, mp[k]})
	}

	sort.Slice(kvPairs, func(i, j int) bool {
		return kvPairs[i].value > kvPairs[j].value
	})

	rankMap := make(map[int]int)
	for i, kv := range kvPairs {
		rankMap[kv.key] = i
	}

	// Create a result slice to store the rankings for the given keys
	result := make([]int, len(keys))
	for i, key := range keys {
		result[i] = rankMap[key] 
	}

	return result
}
