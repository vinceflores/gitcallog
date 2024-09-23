package internal

import "time"



type CalDataPoint struct {
	Date           time.Time
	CommitMessages []string
	// refactor to CommtCount
	CommitCount          float64
}

