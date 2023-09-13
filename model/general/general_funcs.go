package general

import (
	"math/bits"
)

func BitLength(num int) int {
	return bits.Len(uint(num))
}

func ArrContains[T comparable](elems [4]T, value T) bool {
	for _, item := range elems {
		if item == value {
			return true
		}
	}
	return false
}

func Contains[T comparable](elems []T, value T) bool {
	for _, item := range elems {
		if item == value {
			return true
		}
	}
	return false
}
