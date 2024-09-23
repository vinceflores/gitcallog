package internal

var ScaleColors = []string{
	"#161b22", // Less
	"#0e4429",
	"#006d32",
	"#26a641",
	"#39d353", // - More
}

type ViewDataPoint struct {
	actual     float64
	normalized float64
	commits    []string
}
