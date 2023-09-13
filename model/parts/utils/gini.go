package utils

import "math"

func Gini(x []int) float64 {
	total := 0.0
	if len(x) == 0 {
		return 0.0
	}
	for i, xi := range x[:len(x)-1] {
		for _, xj := range x[i+1:] {
			total += math.Abs(float64(xi) - float64(xj))
		}
	}
	avg := Mean(x)
	denom := (math.Pow(float64(len(x)), 2) * avg)
	return total / denom
}

func Mean(x []int) float64 {
	total := 0.0
	for _, xi := range x {
		if xi > 0 {
			total += float64(xi)
		}
	}
	return total / float64(len(x))
}

func Stdev(x []int, mean float64) float64 {
	sum := 0.0
	for _, xi := range x {
		sum += math.Pow(float64(xi)-mean, 2)
	}
	return math.Sqrt(sum / float64(len(x)))
}

func MinAndMax(x []int) (min int, max int) {
	min = x[0]
	max = x[0]
	for _, value := range x {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}
