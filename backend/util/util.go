package util

import (
	"fmt"
)

type Comparable interface {
	LessThan(other Comparable) bool
	GreaterThan(other Comparable) bool
	EqualTo(other Comparable) bool
}

func BinarySearch[T Comparable](arr []T, target T) (T, error) {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := low + (high-low)/2

		if arr[mid].LessThan(target) {
			low = mid + 1
		} else if arr[mid].GreaterThan(target) {
			high = mid - 1
		} else {
			return arr[mid], nil
		}
	}

	return target, fmt.Errorf("not found")
}
