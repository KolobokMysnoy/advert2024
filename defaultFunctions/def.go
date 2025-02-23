package defaultFunc

import "sort"

func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func FindElement(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func QuickSortStart(arr []int) []int {
	sort.Slice(arr, func(i, j int) bool {
		return i > j
	})
	return arr
}
