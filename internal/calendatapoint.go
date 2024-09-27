package internal

import "time"

type CalDataPoint struct {
	Date           time.Time
	CommitMessages []string
	CommitCount    float64
}
