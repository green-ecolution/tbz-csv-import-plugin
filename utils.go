package main

import (
	"math"
)

func CombineToID(a, b float32) int {
	ua := math.Float32bits(a)
	ub := math.Float32bits(b)

	return int((uint(ua) << 32) | uint(ub))
}
