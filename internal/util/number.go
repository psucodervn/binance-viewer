package util

import (
	"math"
	"strconv"
)

const Epsilon = 1e-9

func IsZero(a float64) bool {
	return math.Abs(a) < Epsilon
}

func IsEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		f = math.NaN()
	}
	return f
}

func ParseInt(s string, defaultValue int64) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return v
}
