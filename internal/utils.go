package internal

import "strconv"

func ParseInt(s string) int {
	i , _ :=  strconv.Atoi(s)
	return int(i)
}

func getScaleColor(value float64) string {
	const numColors = 5
	const max = 1.0 // Assume it's normalized between 0.0-1.0
	norm := (value/max)*(numColors-1) 

	if value > 0 && value < 0.5 {
		return ScaleColors[0.5 * (numColors -1) ]
	}

	return ScaleColors[int(norm)]
}