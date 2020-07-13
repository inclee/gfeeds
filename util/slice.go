package util

import (
	"fmt"
)

func IntSliceContain(src []int, item int) bool {
	for _, i := range src {
		if i == item {
			return true
		}
	}
	return false
}

func IntSliceMoveTo(src *[]int, item, toIdx int) (bool, error) {
	if toIdx > len(*src) {
		return false, fmt.Errorf("toIdx it to big than the src length")
	}
	for idx, i := range *src {
		if i == item {
			if idx == toIdx {
				return true, nil
			}
			if idx > toIdx {
				end := append([]int{}, (*src)[toIdx:idx]...)
				end = append(end, (*src)[idx+1:]...)
				*src = append((*src)[:toIdx], item)
				*src = append(*src, end...)
			} else {
				end := append([]int{item}, (*src)[toIdx+1:]...)
				*src = append((*src)[:idx], (*src)[idx+1:toIdx+1]...)
				*src = append(*src, end...)
			}
			return true, nil
		}
	}
	return false, fmt.Errorf("cant find item")
}
