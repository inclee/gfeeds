package util

import (
	"log"
	"testing"
)

func TestSlice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	IntSliceMoveTo(&src, 4, 1)
	log.Println(src)
}
